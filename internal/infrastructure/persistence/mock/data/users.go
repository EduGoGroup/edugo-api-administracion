package data

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// Hash de la contraseña "edugo2024" usando bcrypt
const defaultPasswordHash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

// UUIDs fijos para los usuarios mock
var (
	AdminUserID      = uuid.MustParse("a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01")
	TeacherMathID    = uuid.MustParse("a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a02")
	TeacherScienceID = uuid.MustParse("a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a03")
	Student1ID       = uuid.MustParse("a4eebc99-9c0b-4ef8-bb6d-6bb9bd380a04")
	Student2ID       = uuid.MustParse("a5eebc99-9c0b-4ef8-bb6d-6bb9bd380a05")
	Student3ID       = uuid.MustParse("a6eebc99-9c0b-4ef8-bb6d-6bb9bd380a06")
	Guardian1ID      = uuid.MustParse("a7eebc99-9c0b-4ef8-bb6d-6bb9bd380a07")
	Guardian2ID      = uuid.MustParse("a8eebc99-9c0b-4ef8-bb6d-6bb9bd380a08")
)

// GetUsers retorna un mapa con todos los usuarios mock
func GetUsers() map[uuid.UUID]*entities.User {
	now := time.Now()

	users := map[uuid.UUID]*entities.User{
		AdminUserID: {
			ID:            AdminUserID,
			Email:         "admin@edugo.test",
			PasswordHash:  defaultPasswordHash,
			FirstName:     "Admin",
			LastName:      "Demo",
			Role:          "admin",
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     nil,
		},
		TeacherMathID: {
			ID:            TeacherMathID,
			Email:         "teacher.math@edugo.test",
			PasswordHash:  defaultPasswordHash,
			FirstName:     "María",
			LastName:      "García",
			Role:          "teacher",
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     nil,
		},
		TeacherScienceID: {
			ID:            TeacherScienceID,
			Email:         "teacher.science@edugo.test",
			PasswordHash:  defaultPasswordHash,
			FirstName:     "Juan",
			LastName:      "Pérez",
			Role:          "teacher",
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     nil,
		},
		Student1ID: {
			ID:            Student1ID,
			Email:         "student1@edugo.test",
			PasswordHash:  defaultPasswordHash,
			FirstName:     "Carlos",
			LastName:      "Rodríguez",
			Role:          "student",
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     nil,
		},
		Student2ID: {
			ID:            Student2ID,
			Email:         "student2@edugo.test",
			PasswordHash:  defaultPasswordHash,
			FirstName:     "Ana",
			LastName:      "Martínez",
			Role:          "student",
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     nil,
		},
		Student3ID: {
			ID:            Student3ID,
			Email:         "student3@edugo.test",
			PasswordHash:  defaultPasswordHash,
			FirstName:     "Luis",
			LastName:      "González",
			Role:          "student",
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     nil,
		},
		Guardian1ID: {
			ID:            Guardian1ID,
			Email:         "guardian1@edugo.test",
			PasswordHash:  defaultPasswordHash,
			FirstName:     "Roberto",
			LastName:      "Fernández",
			Role:          "guardian",
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     nil,
		},
		Guardian2ID: {
			ID:            Guardian2ID,
			Email:         "guardian2@edugo.test",
			PasswordHash:  defaultPasswordHash,
			FirstName:     "Patricia",
			LastName:      "López",
			Role:          "guardian",
			IsActive:      true,
			EmailVerified: true,
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     nil,
		},
	}

	return users
}

// GetUserByEmail busca un usuario por email
func GetUserByEmail(email string) *entities.User {
	users := GetUsers()
	for _, user := range users {
		if user.Email == email {
			return user
		}
	}
	return nil
}

// GetUsersByRole retorna todos los usuarios de un rol específico
func GetUsersByRole(role string) []*entities.User {
	users := GetUsers()
	var result []*entities.User

	for _, user := range users {
		if user.Role == role {
			result = append(result, user)
		}
	}

	return result
}
