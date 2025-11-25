// Package handler contiene los handlers HTTP para autenticación
package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/service"
)

// AuthHandler maneja los endpoints de autenticación (login, logout, refresh)
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler crea una nueva instancia de AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary Login de usuario
// @Description Autentica un usuario y retorna tokens JWT
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Credenciales de login"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse "Request inválido"
// @Failure 401 {object} dto.ErrorResponse "Credenciales inválidas"
// @Failure 403 {object} dto.ErrorResponse "Usuario inactivo"
// @Failure 500 {object} dto.ErrorResponse "Error interno"
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Email y password son requeridos",
			Code:    "INVALID_REQUEST",
		})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Credenciales inválidas",
				Code:    "INVALID_CREDENTIALS",
			})
		case errors.Is(err, service.ErrUserInactive):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Usuario inactivo",
				Code:    "USER_INACTIVE",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Error en el proceso de autenticación",
				Code:    "AUTH_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Refresh godoc
// @Summary Refrescar tokens
// @Description Genera nuevos tokens usando el refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse "Request inválido"
// @Failure 401 {object} dto.ErrorResponse "Refresh token inválido"
// @Failure 403 {object} dto.ErrorResponse "Usuario inactivo"
// @Failure 500 {object} dto.ErrorResponse "Error interno"
// @Router /v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Refresh token es requerido",
			Code:    "INVALID_REQUEST",
		})
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Refresh token inválido o expirado",
				Code:    "INVALID_REFRESH_TOKEN",
			})
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "unauthorized",
				Message: "Usuario no encontrado",
				Code:    "USER_NOT_FOUND",
			})
		case errors.Is(err, service.ErrUserInactive):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error:   "forbidden",
				Message: "Usuario inactivo",
				Code:    "USER_INACTIVE",
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   "internal_error",
				Message: "Error refrescando tokens",
				Code:    "REFRESH_ERROR",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary Logout de usuario
// @Description Invalida el access token actual
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]string "Logout exitoso"
// @Failure 400 {object} dto.ErrorResponse "Token no proporcionado"
// @Failure 500 {object} dto.ErrorResponse "Error interno"
// @Router /v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Obtener token del header Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Token de autorización requerido",
			Code:    "TOKEN_REQUIRED",
		})
		return
	}

	// Remover "Bearer " del inicio
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Token vacío",
			Code:    "EMPTY_TOKEN",
		})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Error en logout",
			Code:    "LOGOUT_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout exitoso",
	})
}

// RegisterRoutes registra las rutas del handler de autenticación
func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}
}
