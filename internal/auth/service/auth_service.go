// Package service contiene la lógica de negocio de autenticación
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/repository"
	"github.com/EduGoGroup/edugo-api-administracion/internal/shared/crypto"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/google/uuid"
)

// Errores del servicio de autenticación
var (
	ErrInvalidCredentials  = errors.New("credenciales inválidas")
	ErrUserNotFound        = errors.New("usuario no encontrado")
	ErrUserInactive        = errors.New("usuario inactivo")
	ErrInvalidRefreshToken = errors.New("refresh token inválido")
	ErrNoMembership        = errors.New("no tiene membresía activa en esta escuela")
	ErrInvalidSchoolID     = errors.New("school_id inválido")
)

// AuthService define la interfaz del servicio de autenticación
type AuthService interface {
	// Login valida credenciales y retorna tokens
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)

	// Logout invalida los tokens del usuario
	Logout(ctx context.Context, accessToken string) error

	// SwitchContext cambia el contexto de escuela del usuario
	// Valida que el usuario tenga membresía activa en la escuela destino
	SwitchContext(ctx context.Context, userID, targetSchoolID string) (*dto.SwitchContextResponse, error)

	// RefreshToken genera un nuevo access token usando el refresh token
	// Retorna RefreshResponse (solo access_token) para compatibilidad con api-mobile
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
}

// authService implementa AuthService
type authService struct {
	membershipRepo repository.UnitMembershipRepository
	userRepo       repository.UserRepository
	tokenService   *TokenService
	passwordHasher *crypto.PasswordHasher
	logger         logger.Logger
}

// NewAuthService crea una nueva instancia del servicio
func NewAuthService(
	membershipRepo repository.UnitMembershipRepository,
	userRepo repository.UserRepository,
	tokenService *TokenService,
	passwordHasher *crypto.PasswordHasher,
	logger logger.Logger,
) AuthService {
	return &authService{
		membershipRepo: membershipRepo,
		userRepo:       userRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		logger:         logger,
	}
}

// Login valida credenciales y retorna tokens JWT
func (s *authService) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	// 1. Buscar usuario por email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error buscando usuario: %w", err)
	}
	if user == nil {
		s.logger.Warn("intento de login con email inexistente", "email", email)
		return nil, ErrInvalidCredentials
	}

	// 2. Verificar que el usuario está activo
	if !user.IsActive {
		s.logger.Warn("intento de login con usuario inactivo", "email", email, "user_id", user.ID.String())
		return nil, ErrUserInactive
	}

	// 3. Verificar password
	if err := s.passwordHasher.Compare(password, user.PasswordHash); err != nil {
		s.logger.Warn("password incorrecto", "email", email)
		return nil, ErrInvalidCredentials
	}

	// 4. Obtener school_id del usuario (puede ser nil para super_admin)
	schoolID := ""
	if user.SchoolID != nil {
		schoolID = user.SchoolID.String()
	}

	// 5. Generar tokens (incluyendo school_id)
	tokenResponse, err := s.tokenService.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		user.Role,
		schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("error generando tokens: %w", err)
	}

	// 6. Agregar info del usuario a la respuesta (compatible con api-mobile)
	tokenResponse.User = &dto.UserInfo{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		FullName:  user.FirstName + " " + user.LastName,
		Role:      user.Role,
		SchoolID:  schoolID,
	}

	s.logger.Info("user logged in",
		"entity_type", "auth_session",
		"user_id", user.ID.String(),
		"email", user.Email,
		"role", user.Role,
		"school_id", schoolID,
	)

	// 7. Actualizar último login (fire and forget)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		user.UpdatedAt = time.Now()
		if err := s.userRepo.Update(ctx, user); err != nil {
			s.logger.Warn("error actualizando último login", "error", err)
		}
	}()

	return tokenResponse, nil
}

// Logout invalida el access token agregándolo a la blacklist
func (s *authService) Logout(ctx context.Context, accessToken string) error {
	// Revocar el token (agregarlo a blacklist)
	if err := s.tokenService.RevokeToken(ctx, accessToken); err != nil {
		return fmt.Errorf("error en logout: %w", err)
	}

	s.logger.Info("user logged out",
		"entity_type", "auth_session",
	)
	return nil
}

// RefreshToken valida el refresh token y genera solo un nuevo access token
// Retorna RefreshResponse para compatibilidad con api-mobile y apple-app
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	// 1. Verificar el refresh token
	response, err := s.tokenService.VerifyToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("error verificando token: %w", err)
	}

	if !response.Valid {
		s.logger.Warn("refresh token inválido", "error", response.Error)
		return nil, ErrInvalidRefreshToken
	}

	// 2. Buscar usuario por ID (el refresh token tiene UserID en Subject)
	userID, err := uuid.Parse(response.UserID)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error buscando usuario: %w", err)
	}
	if user == nil {
		s.logger.Warn("usuario no encontrado para refresh", "user_id", response.UserID)
		return nil, ErrUserNotFound
	}

	// 3. Verificar que sigue activo
	if !user.IsActive {
		s.logger.Warn("refresh token de usuario inactivo", "user_id", user.ID.String())
		return nil, ErrUserInactive
	}

	// 4. Obtener school_id del usuario
	schoolID := ""
	if user.SchoolID != nil {
		schoolID = user.SchoolID.String()
	}

	// 5. Generar solo un nuevo access token (NO revocar el refresh token)
	// El refresh token sigue siendo válido hasta su expiración
	refreshResponse, err := s.tokenService.GenerateAccessToken(
		user.ID.String(),
		user.Email,
		user.Role,
		schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("error generando nuevo access token: %w", err)
	}

	s.logger.Info("token refreshed",
		"entity_type", "auth_token",
		"user_id", user.ID.String(),
		"school_id", schoolID,
	)

	return refreshResponse, nil
}

// SwitchContext cambia el contexto de escuela del usuario
// Valida que el usuario tenga una membresía activa en la escuela destino
func (s *authService) SwitchContext(ctx context.Context, userID, targetSchoolID string) (*dto.SwitchContextResponse, error) {
	// 1. Parsear y validar UUIDs
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("user_id inválido: %w", err)
	}

	schoolUUID, err := uuid.Parse(targetSchoolID)
	if err != nil {
		return nil, ErrInvalidSchoolID
	}

	// 2. Verificar que el usuario existe y está activo
	user, err := s.userRepo.FindByID(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error buscando usuario: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// 3. Verificar que el usuario tiene membresía activa en la escuela destino
	membership, err := s.membershipRepo.FindByUserAndSchool(ctx, userUUID, schoolUUID)
	if err != nil {
		return nil, fmt.Errorf("error verificando membresía: %w", err)
	}
	if membership == nil {
		s.logger.Warn("intento de switch-context sin membresía",
			"user_id", userID,
			"target_school_id", targetSchoolID,
		)
		return nil, ErrNoMembership
	}

	// 4. Generar nuevos tokens con el nuevo school_id y rol de la membresía
	tokenResponse, err := s.tokenService.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		membership.Role, // Usar el rol de la membresía en esa escuela
		targetSchoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("error generando tokens: %w", err)
	}

	s.logger.Info("context switched",
		"entity_type", "auth_context",
		"user_id", userID,
		"new_school_id", targetSchoolID,
		"new_role", membership.Role,
	)

	// 5. Construir respuesta
	return &dto.SwitchContextResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
		TokenType:    tokenResponse.TokenType,
		Context: &dto.ContextInfo{
			SchoolID: targetSchoolID,
			Role:     membership.Role,
			UserID:   userID,
			Email:    user.Email,
		},
	}, nil
}
