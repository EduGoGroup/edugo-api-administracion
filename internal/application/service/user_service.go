package service

import (
	"context"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/shared/crypto"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/common/errors"
	"github.com/EduGoGroup/edugo-shared/common/types/enum"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

// UserService define las operaciones de negocio para usuarios
type UserService interface {
	// CreateUser crea un nuevo usuario
	CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error)

	// GetUser obtiene un usuario por ID
	GetUser(ctx context.Context, id string) (*dto.UserResponse, error)

	// GetUserByEmail obtiene un usuario por email
	GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error)

	// UpdateUser actualiza un usuario
	UpdateUser(ctx context.Context, id string, req dto.UpdateUserRequest) (*dto.UserResponse, error)

	// DeleteUser elimina un usuario
	DeleteUser(ctx context.Context, id string) error
}

// userService implementa UserService
type userService struct {
	userRepo       repository.UserRepository
	passwordHasher *crypto.PasswordHasher
	logger         logger.Logger
}

// NewUserService crea un nuevo UserService
func NewUserService(
	userRepo repository.UserRepository,
	logger logger.Logger,
) UserService {
	return &userService{
		userRepo:       userRepo,
		passwordHasher: crypto.NewPasswordHasher(12), // bcrypt cost 12 para producción
		logger:         logger,
	}
}

// CreateUser implementa la creación de un usuario
func (s *userService) CreateUser(
	ctx context.Context,
	req dto.CreateUserRequest,
) (*dto.UserResponse, error) {
	// 1. Validar request
	if err := req.Validate(); err != nil {
		s.logger.Warn("validation failed", "error", err)
		return nil, err
	}

	// 2. Verificar si ya existe un usuario con ese email
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.NewDatabaseError("check user", err)
	}

	if exists {
		return nil, errors.NewAlreadyExistsError("user").
			WithField("email", req.Email)
	}

	// 3. Validar role (lógica de negocio movida del entity)
	role := enum.SystemRole(req.Role)
	if !role.IsValid() {
		return nil, errors.NewValidationError("invalid role").
			WithField("role", req.Role)
	}

	// No permitir crear admin users (regla de negocio)
	if role == enum.SystemRoleAdmin {
		return nil, errors.NewBusinessRuleError("cannot create admin users through this endpoint")
	}

	// 4. Validar y hashear password
	if err := s.passwordHasher.Validate(req.Password); err != nil {
		return nil, errors.NewValidationError(err.Error())
	}

	passwordHash, err := s.passwordHasher.Hash(req.Password)
	if err != nil {
		return nil, errors.NewDatabaseError("hash password", err)
	}

	// 5. Crear entidad de infrastructure
	now := time.Now()
	user := &entities.User{
		ID:            uuid.New(),
		Email:         req.Email,
		PasswordHash:  passwordHash,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          req.Role,
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
		DeletedAt:     nil,
	}

	// 5. Persistir en repositorio
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("create user", err)
	}

	s.logger.Info("entity created",
		"entity_type", "user",
		"entity_id", user.ID.String(),
		"email", user.Email,
		"role", user.Role,
	)

	// 6. Retornar DTO de respuesta
	return dto.ToUserResponse(user), nil
}

// GetUser obtiene un usuario por ID
func (s *userService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id format")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("find user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("user").WithField("id", id)
	}

	return dto.ToUserResponse(user), nil
}

// GetUserByEmail obtiene un usuario por email
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewDatabaseError("find user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("user").WithField("email", email)
	}

	return dto.ToUserResponse(user), nil
}

// UpdateUser actualiza un usuario
func (s *userService) UpdateUser(
	ctx context.Context,
	id string,
	req dto.UpdateUserRequest,
) (*dto.UserResponse, error) {
	// Validar request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Buscar usuario
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id format")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("database error",
			"operation", "find_user",
			"user_id", userID,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError("find user", err)
	}
	if user == nil {
		return nil, errors.NewNotFoundError("user")
	}

	// Actualizar campos (lógica de negocio movida del entity)
	if req.FirstName != nil && req.LastName != nil {
		user.FirstName = *req.FirstName
		user.LastName = *req.LastName
	}

	if req.Role != nil {
		// Validaciones (antes estaban en entity.ChangeRole)
		newRole := enum.SystemRole(*req.Role)
		if !newRole.IsValid() {
			return nil, errors.NewValidationError("invalid role")
		}

		if user.Role == *req.Role {
			return nil, errors.NewBusinessRuleError("new role is the same as current role")
		}

		// No permitir promover a admin
		if newRole == enum.SystemRoleAdmin {
			return nil, errors.NewBusinessRuleError("cannot promote to admin role")
		}

		user.Role = *req.Role
	}

	if req.IsActive != nil {
		// Validaciones (antes estaban en entity.Activate/Deactivate)
		if *req.IsActive && user.IsActive {
			return nil, errors.NewBusinessRuleError("user is already active")
		}
		if !*req.IsActive && !user.IsActive {
			return nil, errors.NewBusinessRuleError("user is already inactive")
		}

		user.IsActive = *req.IsActive
	}

	// Actualizar timestamp
	user.UpdatedAt = time.Now()

	// Persistir
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.NewDatabaseError("update user", err)
	}

	updatedFields := []string{}
	if req.FirstName != nil && req.LastName != nil {
		updatedFields = append(updatedFields, "first_name", "last_name")
	}
	if req.Role != nil {
		updatedFields = append(updatedFields, "role")
	}
	if req.IsActive != nil {
		updatedFields = append(updatedFields, "is_active")
	}

	s.logger.Info("entity updated",
		"entity_type", "user",
		"entity_id", user.ID.String(),
		"fields_updated", updatedFields,
	)

	return dto.ToUserResponse(user), nil
}

// DeleteUser elimina un usuario
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.NewValidationError("invalid user_id format")
	}

	// Verificar que existe
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("database error",
			"operation", "find_user",
			"user_id", userID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError("find user", err)
	}
	if user == nil {
		return errors.NewNotFoundError("user")
	}

	// Soft delete
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return errors.NewDatabaseError("delete user", err)
	}

	s.logger.Info("entity deleted",
		"entity_type", "user",
		"entity_id", userID.String(),
	)

	return nil
}
