package data

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// UUIDs fijos para las relaciones guardian_relations mock
var (
	GuardianRelation1ID = uuid.MustParse("a1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	GuardianRelation2ID = uuid.MustParse("a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	GuardianRelation3ID = uuid.MustParse("a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")
)

// GetGuardianRelations retorna un mapa con todas las relaciones guardian-student mock
func GetGuardianRelations() map[uuid.UUID]*entities.GuardianRelation {
	now := time.Now()

	relations := map[uuid.UUID]*entities.GuardianRelation{
		// Roberto Fernández (Guardian1) → Carlos Rodríguez (Student1), Type: father
		GuardianRelation1ID: {
			ID:               GuardianRelation1ID,
			GuardianID:       Guardian1ID, // Roberto Fernández (a7..07)
			StudentID:        Student1ID,  // Carlos Rodríguez (a4..04)
			RelationshipType: "father",
			IsActive:         true,
			CreatedAt:        now,
			UpdatedAt:        now,
			CreatedBy:        "system",
		},
		// Patricia López (Guardian2) → Ana Martínez (Student2), Type: mother
		GuardianRelation2ID: {
			ID:               GuardianRelation2ID,
			GuardianID:       Guardian2ID, // Patricia López (a8..08)
			StudentID:        Student2ID,  // Ana Martínez (a5..05)
			RelationshipType: "mother",
			IsActive:         true,
			CreatedAt:        now,
			UpdatedAt:        now,
			CreatedBy:        "system",
		},
		// Roberto Fernández (Guardian1) → Luis González (Student3), Type: legal_guardian
		GuardianRelation3ID: {
			ID:               GuardianRelation3ID,
			GuardianID:       Guardian1ID, // Roberto Fernández (a7..07)
			StudentID:        Student3ID,  // Luis González (a6..06)
			RelationshipType: "legal_guardian",
			IsActive:         true,
			CreatedAt:        now,
			UpdatedAt:        now,
			CreatedBy:        "system",
		},
	}

	return relations
}

// GetRelationsByGuardian retorna todas las relaciones de un guardian específico
func GetRelationsByGuardian(guardianID uuid.UUID) []*entities.GuardianRelation {
	relations := GetGuardianRelations()
	var result []*entities.GuardianRelation

	for _, relation := range relations {
		if relation.GuardianID == guardianID && relation.IsActive {
			result = append(result, relation)
		}
	}

	return result
}

// GetRelationsByStudent retorna todas las relaciones de un estudiante específico
func GetRelationsByStudent(studentID uuid.UUID) []*entities.GuardianRelation {
	relations := GetGuardianRelations()
	var result []*entities.GuardianRelation

	for _, relation := range relations {
		if relation.StudentID == studentID && relation.IsActive {
			result = append(result, relation)
		}
	}

	return result
}
