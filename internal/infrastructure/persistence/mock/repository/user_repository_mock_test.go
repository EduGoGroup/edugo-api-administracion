package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestMockUserRepository_FindByEmail(t *testing.T) {
	repo := NewMockUserRepository()
	ctx := context.Background()

	// Test: Buscar usuario admin
	user, err := repo.FindByEmail(ctx, "admin@edugo.test")
	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "admin@edugo.test", user.Email)
	assert.Equal(t, "Admin", user.FirstName)
	assert.Equal(t, "admin", user.Role)

	// Test: Verificar hash de contraseña
	password := "edugo2024"
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	assert.NoError(t, err, "El hash de contraseña debería ser válido para 'edugo2024'")
}

func TestMockUserRepository_ExistsByEmail(t *testing.T) {
	repo := NewMockUserRepository()
	ctx := context.Background()

	exists, err := repo.ExistsByEmail(ctx, "admin@edugo.test")
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByEmail(ctx, "noexiste@edugo.test")
	require.NoError(t, err)
	assert.False(t, exists)
}
