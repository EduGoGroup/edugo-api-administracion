package service

import (
	"context"

	"github.com/EduGoGroup/edugo-api-administracion/internal/application/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/logger"
	"github.com/EduGoGroup/edugo-shared/pkg/types/enum"
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
	userRepo repository.UserRepository
	logger   logger.Logger
}

// NewUserService crea un nuevo UserService
func NewUserService(
	userRepo repository.UserRepository,
	logger logger.Logger,
) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
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

	// 2. Crear value object Email
	email, err := valueobject.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	// 3. Verificar si ya existe un usuario con ese email
	exists, err := s.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to check existing user",
			"error", err,
			"email", email.String(),
		)
		return nil, errors.NewDatabaseError("check user", err)
	}

	if exists {
		return nil, errors.NewAlreadyExistsError("user").
			WithField("email", email.String())
	}

	// 4. Convertir role a enum
	role := enum.SystemRole(req.Role)
	if !role.IsValid() {
		return nil, errors.NewValidationError("invalid role").
			WithField("role", req.Role)
	}

	// 5. Crear entidad de dominio
	user, err := entity.NewUser(
		email,
		req.FirstName,
		req.LastName,
		role,
	)
	if err != nil {
		s.logger.Warn("failed to create user entity", "error", err)
		return nil, err
	}

	// 6. Persistir en repositorio
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("failed to save user",
			"error", err,
			"email", email.String(),
		)
		return nil, errors.NewDatabaseError("create user", err)
	}

	s.logger.Info("user created",
		"user_id", user.ID().String(),
		"email", email.String(),
		"role", role.String(),
	)

	// 7. Retornar DTO de respuesta
	return dto.ToUserResponse(user), nil
}

// GetUser obtiene un usuario por ID
func (s *userService) GetUser(ctx context.Context, id string) (*dto.UserResponse, error) {
	userID, err := valueobject.UserIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id format")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to find user", "error", err, "id", id)
		return nil, errors.NewDatabaseError("find user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("user").WithField("id", id)
	}

	return dto.ToUserResponse(user), nil
}

// GetUserByEmail obtiene un usuario por email
func (s *userService) GetUserByEmail(ctx context.Context, emailStr string) (*dto.UserResponse, error) {
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to find user", "error", err, "email", emailStr)
		return nil, errors.NewDatabaseError("find user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("user").WithField("email", emailStr)
	}

	return dto.ToUserResponse(user), nil
}

// UpdateUser actualiza un usuario (implementación simplificada)
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
	userID, err := valueobject.UserIDFromString(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid user_id format")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.NewNotFoundError("user")
	}

	// Actualizar campos
	if req.FirstName != nil && req.LastName != nil {
		if err := user.UpdateName(*req.FirstName, *req.LastName); err != nil {
			return nil, err
		}
	}

	if req.Role != nil {
		role := enum.SystemRole(*req.Role)
		if err := user.ChangeRole(role); err != nil {
			return nil, err
		}
	}

	if req.IsActive != nil {
		if *req.IsActive {
			if err := user.Activate(); err != nil {
				return nil, err
			}
		} else {
			if err := user.Deactivate(); err != nil {
				return nil, err
			}
		}
	}

	// Persistir
	if err := s.userRepo.Update(ctx, user); err != nil {
		s.logger.Error("failed to update user", "error", err)
		return nil, errors.NewDatabaseError("update user", err)
	}

	s.logger.Info("user updated", "user_id", user.ID().String())

	return dto.ToUserResponse(user), nil
}

// DeleteUser elimina un usuario
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	userID, err := valueobject.UserIDFromString(id)
	if err != nil {
		return errors.NewValidationError("invalid user_id format")
	}

	// Verificar que existe
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return errors.NewNotFoundError("user")
	}

	// Soft delete
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		s.logger.Error("failed to delete user", "error", err, "id", id)
		return errors.NewDatabaseError("delete user", err)
	}

	s.logger.Info("user deleted", "user_id", userID.String())

	return nil
}
