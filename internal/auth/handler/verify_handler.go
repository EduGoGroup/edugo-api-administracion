// Package handler contiene los handlers HTTP para autenticación
package handler

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/dto"
	"github.com/EduGoGroup/edugo-api-administracion/internal/auth/service"
)

// VerifyHandler maneja las solicitudes de verificación de tokens
type VerifyHandler struct {
	tokenService  *service.TokenService
	internalNets  []*net.IPNet
	apiKeys       map[string]string
}

// NewVerifyHandler crea una nueva instancia del handler
func NewVerifyHandler(
	tokenService *service.TokenService,
	internalIPRanges []string,
	apiKeys map[string]string,
) *VerifyHandler {
	// Parsear rangos CIDR
	var nets []*net.IPNet
	for _, cidr := range internalIPRanges {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			// Si no es CIDR, intentar como IP simple
			ip := net.ParseIP(cidr)
			if ip != nil {
				// Convertir IP a CIDR /32 o /128
				if ip.To4() != nil {
					_, ipNet, _ = net.ParseCIDR(cidr + "/32")
				} else {
					_, ipNet, _ = net.ParseCIDR(cidr + "/128")
				}
			}
		}
		if ipNet != nil {
			nets = append(nets, ipNet)
		}
	}

	return &VerifyHandler{
		tokenService:  tokenService,
		internalNets:  nets,
		apiKeys:       apiKeys,
	}
}

// VerifyToken godoc
// @Summary Verificar token JWT
// @Description Verifica la validez de un token JWT y retorna información del usuario
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyTokenRequest true "Token a verificar"
// @Param X-Service-API-Key header string false "API Key del servicio (para rate limiting diferenciado)"
// @Success 200 {object} dto.VerifyTokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 429 {object} dto.ErrorResponse "Rate limit excedido"
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/verify [post]
func (h *VerifyHandler) VerifyToken(c *gin.Context) {
	startTime := time.Now()

	// 1. Parsear request
	var req dto.VerifyTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Token es requerido",
			Code:    "INVALID_REQUEST",
		})
		return
	}

	// 2. Limpiar token (remover "Bearer " si está presente)
	token := strings.TrimPrefix(req.Token, "Bearer ")
	token = strings.TrimSpace(token)

	if token == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Token vacío",
			Code:    "EMPTY_TOKEN",
		})
		return
	}

	// 3. Verificar token
	response, err := h.tokenService.VerifyToken(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Error verificando token",
			Code:    "VERIFICATION_ERROR",
		})
		return
	}

	// 4. Agregar métricas (header con tiempo de respuesta)
	duration := time.Since(startTime)
	c.Header("X-Response-Time", duration.String())

	// 5. Retornar respuesta
	// Siempre retornar 200, el campo "valid" indica si el token es válido
	c.JSON(http.StatusOK, response)
}

// VerifyTokenBulk godoc
// @Summary Verificar múltiples tokens JWT
// @Description Verifica la validez de múltiples tokens JWT en una sola llamada
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyTokenBulkRequest true "Tokens a verificar"
// @Param X-Service-API-Key header string true "API Key del servicio (requerido para bulk)"
// @Success 200 {object} dto.VerifyTokenBulkResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse "API Key inválida"
// @Failure 429 {object} dto.ErrorResponse "Rate limit excedido"
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/auth/verify-bulk [post]
func (h *VerifyHandler) VerifyTokenBulk(c *gin.Context) {
	// 1. Verificar que es un servicio interno (API Key requerida para bulk)
	if !h.isInternalService(c) {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "API Key requerida para verificación en lote",
			Code:    "API_KEY_REQUIRED",
		})
		return
	}

	// 2. Parsear request
	var req dto.VerifyTokenBulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Lista de tokens es requerida (máximo 100)",
			Code:    "INVALID_REQUEST",
		})
		return
	}

	// 3. Validar cantidad
	if len(req.Tokens) == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Debe proporcionar al menos un token",
			Code:    "EMPTY_TOKENS",
		})
		return
	}

	if len(req.Tokens) > 100 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "bad_request",
			Message: "Máximo 100 tokens por request",
			Code:    "TOO_MANY_TOKENS",
		})
		return
	}

	// 4. Verificar tokens
	response, err := h.tokenService.VerifyTokenBulk(c.Request.Context(), req.Tokens)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Error verificando tokens",
			Code:    "VERIFICATION_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// isInternalService verifica si el request viene de un servicio interno
func (h *VerifyHandler) isInternalService(c *gin.Context) bool {
	// Verificar API Key
	apiKey := c.GetHeader("X-Service-API-Key")
	if apiKey != "" {
		for _, key := range h.apiKeys {
			if key == apiKey {
				return true
			}
		}
	}

	// Verificar IP
	clientIP := net.ParseIP(c.ClientIP())
	if clientIP != nil {
		for _, ipNet := range h.internalNets {
			if ipNet.Contains(clientIP) {
				return true
			}
		}
	}

	return false
}

// IsInternalService expone la verificación para uso en middlewares
func (h *VerifyHandler) IsInternalService(c *gin.Context) bool {
	return h.isInternalService(c)
}

// RegisterRoutes registra las rutas del handler
func (h *VerifyHandler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/verify", h.VerifyToken)
		auth.POST("/verify-bulk", h.VerifyTokenBulk)
	}
}
