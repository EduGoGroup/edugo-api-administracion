package data

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	// UUIDs fijos para las memberships mock
	MembershipTeacherMaria  = uuid.MustParse("d1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	MembershipTeacherJuan   = uuid.MustParse("d2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	MembershipStudentCarlos = uuid.MustParse("d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")
	MembershipStudentAna    = uuid.MustParse("d4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44")
	MembershipStudentLuis   = uuid.MustParse("d5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55")

	// Nota: Se usan las siguientes variables ya definidas en otros archivos del paquete:
	// - users.go: TeacherMathID, TeacherScienceID, Student1ID, Student2ID, Student3ID
	// - academic_units.go: schoolPrimariaID, seccionPrimer1A, seccionPrimer1B
)

// GetMemberships retorna un mapa con todas las memberships mock
func GetMemberships() map[uuid.UUID]*entities.Membership {
	now := time.Now()
	enrolledDate := time.Date(2024, 3, 1, 9, 0, 0, 0, time.UTC)

	// Helper para UUIDs opcionales
	uuidPtr := func(id uuid.UUID) *uuid.UUID { return &id }

	memberships := map[uuid.UUID]*entities.Membership{
		// Teacher María García → Escuela Primaria → Sección A
		MembershipTeacherMaria: {
			ID:             MembershipTeacherMaria,
			UserID:         TeacherMathID,            // a2..02
			SchoolID:       schoolPrimariaID,         // b1..11
			AcademicUnitID: uuidPtr(seccionPrimer1A), // c4..44
			Role:           "teacher",
			IsActive:       true,
			EnrolledAt:     enrolledDate,
			WithdrawnAt:    nil,
			Metadata:       []byte(`{"subject": "Matemáticas", "schedule": "Lunes a Viernes 8:00-12:00", "classroom": "101"}`),
			CreatedAt:      now,
			UpdatedAt:      now,
		},

		// Teacher Juan Pérez → Escuela Primaria → Sección B
		MembershipTeacherJuan: {
			ID:             MembershipTeacherJuan,
			UserID:         TeacherScienceID,         // a3..03
			SchoolID:       schoolPrimariaID,         // b1..11
			AcademicUnitID: uuidPtr(seccionPrimer1B), // c5..55
			Role:           "teacher",
			IsActive:       true,
			EnrolledAt:     enrolledDate,
			WithdrawnAt:    nil,
			Metadata:       []byte(`{"subject": "Ciencias Naturales", "schedule": "Lunes a Viernes 8:00-12:00", "classroom": "102"}`),
			CreatedAt:      now,
			UpdatedAt:      now,
		},

		// Student Carlos → Escuela Primaria → Sección A
		MembershipStudentCarlos: {
			ID:             MembershipStudentCarlos,
			UserID:         Student1ID,               // a4..04
			SchoolID:       schoolPrimariaID,         // b1..11
			AcademicUnitID: uuidPtr(seccionPrimer1A), // c4..44
			Role:           "student",
			IsActive:       true,
			EnrolledAt:     enrolledDate,
			WithdrawnAt:    nil,
			Metadata:       []byte(`{"enrollment_number": "2024-P1A-001", "birth_date": "2017-05-15", "guardian_contact": "+54-11-1111-1111"}`),
			CreatedAt:      now,
			UpdatedAt:      now,
		},

		// Student Ana → Escuela Primaria → Sección A
		MembershipStudentAna: {
			ID:             MembershipStudentAna,
			UserID:         Student2ID,               // a5..05
			SchoolID:       schoolPrimariaID,         // b1..11
			AcademicUnitID: uuidPtr(seccionPrimer1A), // c4..44
			Role:           "student",
			IsActive:       true,
			EnrolledAt:     enrolledDate,
			WithdrawnAt:    nil,
			Metadata:       []byte(`{"enrollment_number": "2024-P1A-002", "birth_date": "2017-08-22", "guardian_contact": "+54-11-2222-2222"}`),
			CreatedAt:      now,
			UpdatedAt:      now,
		},

		// Student Luis → Escuela Primaria → Sección B
		MembershipStudentLuis: {
			ID:             MembershipStudentLuis,
			UserID:         Student3ID,               // a6..06
			SchoolID:       schoolPrimariaID,         // b1..11
			AcademicUnitID: uuidPtr(seccionPrimer1B), // c5..55
			Role:           "student",
			IsActive:       true,
			EnrolledAt:     enrolledDate,
			WithdrawnAt:    nil,
			Metadata:       []byte(`{"enrollment_number": "2024-P1B-001", "birth_date": "2017-03-10", "guardian_contact": "+54-11-3333-3333"}`),
			CreatedAt:      now,
			UpdatedAt:      now,
		},
	}

	return memberships
}

// GetMembershipsByUser retorna todas las memberships de un usuario específico
func GetMembershipsByUser(userID uuid.UUID) []*entities.Membership {
	memberships := GetMemberships()
	var result []*entities.Membership

	for _, membership := range memberships {
		if membership.UserID == userID {
			result = append(result, membership)
		}
	}

	return result
}

// GetMembershipsByUnit retorna todas las memberships de una unidad académica específica
func GetMembershipsByUnit(unitID uuid.UUID) []*entities.Membership {
	memberships := GetMemberships()
	var result []*entities.Membership

	for _, membership := range memberships {
		if membership.AcademicUnitID != nil && *membership.AcademicUnitID == unitID {
			result = append(result, membership)
		}
	}

	return result
}

// GetMembershipsByRole retorna todas las memberships de un rol específico en una escuela
func GetMembershipsByRole(schoolID uuid.UUID, role string) []*entities.Membership {
	memberships := GetMemberships()
	var result []*entities.Membership

	for _, membership := range memberships {
		if membership.SchoolID == schoolID && membership.Role == role {
			result = append(result, membership)
		}
	}

	return result
}

// GetMembershipByID retorna una membership por su ID
func GetMembershipByID(id uuid.UUID) *entities.Membership {
	memberships := GetMemberships()
	return memberships[id]
}

// GetActiveMemberships retorna solo las memberships activas
func GetActiveMemberships() []*entities.Membership {
	memberships := GetMemberships()
	var result []*entities.Membership

	for _, membership := range memberships {
		if membership.IsActive && membership.WithdrawnAt == nil {
			result = append(result, membership)
		}
	}

	return result
}

// GetMembershipsBySchool retorna todas las memberships de una escuela específica
func GetMembershipsBySchool(schoolID uuid.UUID) []*entities.Membership {
	memberships := GetMemberships()
	var result []*entities.Membership

	for _, membership := range memberships {
		if membership.SchoolID == schoolID {
			result = append(result, membership)
		}
	}

	return result
}
