package data

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	// UUIDs fijos para las escuelas de prueba
	schoolID1 = uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	schoolID2 = uuid.MustParse("b2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	schoolID3 = uuid.MustParse("b3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")

	// Datos comunes
	baseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	// Helpers para strings opcionales
	strPtr = func(s string) *string { return &s }
)

// GetSchools retorna un map con las 3 escuelas de prueba
func GetSchools() map[uuid.UUID]*entities.School {
	schools := make(map[uuid.UUID]*entities.School)

	// Escuela Primaria Demo
	schools[schoolID1] = &entities.School{
		ID:               schoolID1,
		Name:             "Escuela Primaria Demo",
		Code:             "SCH_PRI_001",
		Address:          strPtr("Calle Principal 123"),
		City:             strPtr("Buenos Aires"),
		Country:          "Argentina",
		Phone:            strPtr("+54-11-1234-5678"),
		Email:            strPtr("contacto@primaria.test"),
		Metadata:         []byte(`{}`),
		IsActive:         true,
		SubscriptionTier: "basic",
		MaxTeachers:      20,
		MaxStudents:      200,
		CreatedAt:        baseTime,
		UpdatedAt:        baseTime,
		DeletedAt:        nil,
	}

	// Colegio Secundario Demo
	schools[schoolID2] = &entities.School{
		ID:               schoolID2,
		Name:             "Colegio Secundario Demo",
		Code:             "SCH_SEC_001",
		Address:          strPtr("Avenida Libertador 456"),
		City:             strPtr("Buenos Aires"),
		Country:          "Argentina",
		Phone:            strPtr("+54-11-8765-4321"),
		Email:            strPtr("info@secundario.test"),
		Metadata:         []byte(`{}`),
		IsActive:         true,
		SubscriptionTier: "premium",
		MaxTeachers:      50,
		MaxStudents:      500,
		CreatedAt:        baseTime,
		UpdatedAt:        baseTime,
		DeletedAt:        nil,
	}

	// Instituto Técnico Demo
	schools[schoolID3] = &entities.School{
		ID:               schoolID3,
		Name:             "Instituto Técnico Demo",
		Code:             "SCH_TEC_001",
		Address:          strPtr("Boulevard Tecnológico 789"),
		City:             strPtr("Córdoba"),
		Country:          "Argentina",
		Phone:            strPtr("+54-351-999-8888"),
		Email:            strPtr("admin@tecnico.test"),
		Metadata:         []byte(`{}`),
		IsActive:         true,
		SubscriptionTier: "premium",
		MaxTeachers:      50,
		MaxStudents:      500,
		CreatedAt:        baseTime,
		UpdatedAt:        baseTime,
		DeletedAt:        nil,
	}

	return schools
}

// GetSchoolByCode retorna una escuela por su código único
func GetSchoolByCode(code string) *entities.School {
	schools := GetSchools()
	for _, school := range schools {
		if school.Code == code {
			return school
		}
	}
	return nil
}

// GetSchoolByName retorna una escuela por su nombre
func GetSchoolByName(name string) *entities.School {
	schools := GetSchools()
	for _, school := range schools {
		if school.Name == name {
			return school
		}
	}
	return nil
}

// GetSchoolByID retorna una escuela por su ID
func GetSchoolByID(id uuid.UUID) *entities.School {
	schools := GetSchools()
	return schools[id]
}
