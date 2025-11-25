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
)

// AuthService define la interfaz del servicio de autenticación
type AuthService interface {
	// Login valida credenciales y retorna tokens
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)

	// Logout invalida los tokens del usuario
	Logout(ctx context.Context, accessToken string) error

	// RefreshToken genera nuevos tokens usando el refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
}

// authService implementa AuthService
type authService struct {
	userRepo       repository.UserRepository
	tokenService   *TokenService
	passwordHasher *crypto.PasswordHasher
	logger         logger.Logger
}

// NewAuthService crea una nueva instancia del servicio
func NewAuthService(
	userRepo repository.UserRepository,
	tokenService *TokenService,
	passwordHasher *crypto.PasswordHasher,
	logger logger.Logger,
) AuthService {
	return &authService{
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
		s.logger.Error("error buscando usuario", "error", err, "email", email)
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

	// 4. Generar tokens
	tokenResponse, err := s.tokenService.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		user.Role,
	)
	if err != nil {
		s.logger.Error("error generando tokens", "error", err, "user_id", user.ID.String())
		return nil, fmt.Errorf("error generando tokens: %w", err)
	}

	// 5. Agregar info del usuario a la respuesta
	tokenResponse.User = &dto.UserInfo{
		ID:            user.ID.String(),
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Role:          user.Role,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
	}

	s.logger.Info("login exitoso",
		"user_id", user.ID.String(),
		"email", email,
		"role", user.Role,
	)

	// 6. Actualizar último login (fire and forget)
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
		s.logger.Error("error revocando token", "error", err)
		return fmt.Errorf("error en logout: %w", err)
	}

	s.logger.Info("logout exitoso")
	return nil
}

// RefreshToken valida el refresh token y genera nuevos tokens
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	// 1. Verificar el refresh token
	response, err := s.tokenService.VerifyToken(ctx, refreshToken)
	if err != nil {
		s.logger.Error("error verificando refresh token", "error", err)
		return nil, fmt.Errorf("error verificando token: %w", err)
	}

	if !response.Valid {
		s.logger.Warn("refresh token inválido", "error", response.Error)
		return nil, ErrInvalidRefreshToken
	}

	// 2. Buscar usuario por ID (el refresh token tiene UserID en Subject)
	userID, err := uuid.Parse(response.UserID)
	if err != nil {
		s.logger.Error("error parseando user_id del token", "error", err, "user_id", response.UserID)
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Error("error buscando usuario para refresh", "error", err)
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

	// 4. Revocar el refresh token anterior
	_ = s.tokenService.RevokeToken(ctx, refreshToken)

	// 5. Generar nuevos tokens
	tokenResponse, err := s.tokenService.GenerateTokenPair(
		user.ID.String(),
		user.Email,
		user.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("error generando nuevos tokens: %w", err)
	}

	// 6. Agregar info del usuario
	tokenResponse.User = &dto.UserInfo{
		ID:            user.ID.String(),
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Role:          user.Role,
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
	}

	s.logger.Info("token refresh exitoso", "user_id", user.ID.String())

	return tokenResponse, nil
}
