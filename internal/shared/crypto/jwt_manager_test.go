package crypto

import (
	"strings"
	"testing"
	"time"
)

func TestNewJWTManager(t *testing.T) {
	tests := []struct {
		name        string
		config      JWTConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "configuración válida",
			config: JWTConfig{
				Secret:               "test-secret-key-minimum-32-characters-long",
				Issuer:               "edugo-central",
				AccessTokenDuration:  15 * time.Minute,
				RefreshTokenDuration: 7 * 24 * time.Hour,
			},
			expectError: false,
		},
		{
			name: "secret muy corto",
			config: JWTConfig{
				Secret: "short",
				Issuer: "edugo-central",
			},
			expectError: true,
			errorMsg:    "32 caracteres",
		},
		{
			name: "issuer vacío",
			config: JWTConfig{
				Secret: "test-secret-key-minimum-32-characters-long",
				Issuer: "",
			},
			expectError: true,
			errorMsg:    "issuer es requerido",
		},
		{
			name: "duraciones por defecto",
			config: JWTConfig{
				Secret: "test-secret-key-minimum-32-characters-long",
				Issuer: "edugo-central",
				// AccessTokenDuration y RefreshTokenDuration = 0
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewJWTManager(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("esperaba error pero no lo hubo")
				}
				if tt.errorMsg != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errorMsg) {
						t.Errorf("mensaje de error incorrecto: %v", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("no esperaba error: %v", err)
				}
				if manager == nil {
					t.Error("manager es nil")
				}
			}
		})
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	manager := createTestManager(t)

	// Generar token
	userID := "user-123"
	email := "test@edugo.com"
	role := "teacher"

	token, expiresAt, err := manager.GenerateAccessToken(userID, email, role, "")
	if err != nil {
		t.Fatalf("error generando token: %v", err)
	}

	if token == "" {
		t.Error("token está vacío")
	}

	if expiresAt.Before(time.Now()) {
		t.Error("expiración debe ser en el futuro")
	}

	// Validar token
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("error validando token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID incorrecto: esperado %s, obtenido %s", userID, claims.UserID)
	}

	if claims.Email != email {
		t.Errorf("Email incorrecto: esperado %s, obtenido %s", email, claims.Email)
	}

	if claims.Role != role {
		t.Errorf("Role incorrecto: esperado %s, obtenido %s", role, claims.Role)
	}

	if claims.Issuer != "edugo-central" {
		t.Errorf("Issuer incorrecto: esperado edugo-central, obtenido %s", claims.Issuer)
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := createTestManager(t)

	userID := "user-456"
	token, expiresAt, err := manager.GenerateRefreshToken(userID)

	if err != nil {
		t.Fatalf("error generando refresh token: %v", err)
	}

	if token == "" {
		t.Error("refresh token está vacío")
	}

	// Refresh token debe expirar después del access token (7 días por defecto)
	expectedExpiry := time.Now().Add(7 * 24 * time.Hour)
	if expiresAt.Before(expectedExpiry.Add(-1 * time.Minute)) {
		t.Error("refresh token expira demasiado pronto")
	}

	// Validar que el token es válido
	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatalf("error validando refresh token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID incorrecto en refresh token")
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	manager := createTestManager(t)

	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "token vacío",
			token: "",
		},
		{
			name:  "token malformado",
			token: "not.a.valid.token",
		},
		{
			name:  "token con firma incorrecta",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.ValidateToken(tt.token)
			if err == nil {
				t.Error("esperaba error")
			}
		})
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Crear manager con duración muy corta
	config := JWTConfig{
		Secret:              "test-secret-key-minimum-32-characters-long",
		Issuer:              "edugo-central",
		AccessTokenDuration: 1 * time.Millisecond,
	}

	manager, err := NewJWTManager(config)
	if err != nil {
		t.Fatalf("error creando manager: %v", err)
	}

	// Generar token
	token, _, err := manager.GenerateAccessToken("user-1", "test@test.com", "student", "")
	if err != nil {
		t.Fatalf("error generando token: %v", err)
	}

	// Esperar a que expire
	time.Sleep(10 * time.Millisecond)

	// Validar - debe fallar por expirado
	_, err = manager.ValidateToken(token)
	if err != ErrTokenExpired {
		t.Errorf("esperaba ErrTokenExpired, obtuvo: %v", err)
	}
}

func TestValidateToken_WrongIssuer(t *testing.T) {
	// Manager 1 con issuer "edugo-central"
	manager1, _ := NewJWTManager(JWTConfig{
		Secret:              "test-secret-key-minimum-32-characters-long",
		Issuer:              "edugo-central",
		AccessTokenDuration: 1 * time.Hour,
	})

	// Manager 2 con issuer diferente
	manager2, _ := NewJWTManager(JWTConfig{
		Secret:              "test-secret-key-minimum-32-characters-long",
		Issuer:              "otro-issuer",
		AccessTokenDuration: 1 * time.Hour,
	})

	// Generar token con manager2
	token, _, _ := manager2.GenerateAccessToken("user-1", "test@test.com", "student", "")

	// Intentar validar con manager1
	_, err := manager1.ValidateToken(token)
	if err == nil {
		t.Error("esperaba error por issuer incorrecto")
	}

	if !strings.Contains(err.Error(), "issuer") {
		t.Errorf("error debe mencionar issuer: %v", err)
	}
}

func TestGetTokenID(t *testing.T) {
	manager := createTestManager(t)

	token, _, _ := manager.GenerateAccessToken("user-1", "test@test.com", "student", "")

	tokenID, err := manager.GetTokenID(token)
	if err != nil {
		t.Fatalf("error obteniendo token ID: %v", err)
	}

	if tokenID == "" {
		t.Error("token ID está vacío")
	}

	// Verificar formato UUID (36 caracteres con guiones)
	if len(tokenID) != 36 {
		t.Errorf("token ID no parece ser UUID: %s (len=%d)", tokenID, len(tokenID))
	}
}

func TestGetExpirationTime(t *testing.T) {
	manager := createTestManager(t)

	token, expectedExpiry, _ := manager.GenerateAccessToken("user-1", "test@test.com", "student", "")

	expiry, err := manager.GetExpirationTime(token)
	if err != nil {
		t.Fatalf("error obteniendo expiración: %v", err)
	}

	// Permitir 1 segundo de diferencia por tiempo de procesamiento
	diff := expiry.Sub(expectedExpiry)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("expiración incorrecta: esperado %v, obtenido %v", expectedExpiry, expiry)
	}
}

func TestGetConfig(t *testing.T) {
	config := JWTConfig{
		Secret:               "test-secret-key-minimum-32-characters-long",
		Issuer:               "edugo-central",
		AccessTokenDuration:  30 * time.Minute,
		RefreshTokenDuration: 14 * 24 * time.Hour,
	}

	manager, _ := NewJWTManager(config)
	returnedConfig := manager.GetConfig()

	if returnedConfig.Issuer != config.Issuer {
		t.Errorf("issuer incorrecto: esperado %s, obtenido %s", config.Issuer, returnedConfig.Issuer)
	}

	if returnedConfig.AccessTokenDuration != config.AccessTokenDuration {
		t.Errorf("access token duration incorrecta")
	}
}

// Helper function
func createTestManager(t *testing.T) *JWTManager {
	t.Helper()

	manager, err := NewJWTManager(JWTConfig{
		Secret:               "test-secret-key-minimum-32-characters-long",
		Issuer:               "edugo-central",
		AccessTokenDuration:  15 * time.Minute,
		RefreshTokenDuration: 7 * 24 * time.Hour,
	})

	if err != nil {
		t.Fatalf("error creando manager: %v", err)
	}

	return manager
}
