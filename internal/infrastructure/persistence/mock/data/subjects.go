package data

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	// UUIDs fijos para las materias de prueba
	MathematicsID       = uuid.MustParse("e1eebc99-9c0b-4ef8-bb6d-6bb9bd380e11")
	NaturalSciencesID   = uuid.MustParse("e2eebc99-9c0b-4ef8-bb6d-6bb9bd380e22")
	SpanishLiteratureID = uuid.MustParse("e3eebc99-9c0b-4ef8-bb6d-6bb9bd380e33")
	HistoryID           = uuid.MustParse("e4eebc99-9c0b-4ef8-bb6d-6bb9bd380e44")
	ProgrammingID       = uuid.MustParse("e5eebc99-9c0b-4ef8-bb6d-6bb9bd380e55")
	PhysicalEdID        = uuid.MustParse("e6eebc99-9c0b-4ef8-bb6d-6bb9bd380e66")

	// Timestamp base para consistencia en datos mock
	subjectsBaseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

// GetSubjects retorna un mapa con las 6 materias de prueba
func GetSubjects() map[uuid.UUID]*entities.Subject {
	subjects := make(map[uuid.UUID]*entities.Subject)

	// Helper para crear pointers a strings
	sp := func(s string) *string { return &s }

	// Matemáticas
	subjects[MathematicsID] = &entities.Subject{
		ID:          MathematicsID,
		Name:        "Matemáticas",
		Description: sp("Materia de ciencias exactas que cubre álgebra, geometría, cálculo y trigonometría"),
		Metadata:    sp(`{"area":"Ciencias Exactas","level":"Secundario","difficulty":"Intermedio"}`),
		IsActive:    true,
		CreatedAt:   subjectsBaseTime,
		UpdatedAt:   subjectsBaseTime,
	}

	// Ciencias Naturales
	subjects[NaturalSciencesID] = &entities.Subject{
		ID:          NaturalSciencesID,
		Name:        "Ciencias Naturales",
		Description: sp("Estudio de biología, química y física en el contexto de la naturaleza"),
		Metadata:    sp(`{"area":"Ciencias Naturales","level":"Secundario","difficulty":"Intermedio"}`),
		IsActive:    true,
		CreatedAt:   subjectsBaseTime,
		UpdatedAt:   subjectsBaseTime,
	}

	// Lengua y Literatura
	subjects[SpanishLiteratureID] = &entities.Subject{
		ID:          SpanishLiteratureID,
		Name:        "Lengua y Literatura",
		Description: sp("Desarrollo de competencias de lectura, escritura y análisis literario en español"),
		Metadata:    sp(`{"area":"Humanidades","level":"Secundario","difficulty":"Intermedio"}`),
		IsActive:    true,
		CreatedAt:   subjectsBaseTime,
		UpdatedAt:   subjectsBaseTime,
	}

	// Historia
	subjects[HistoryID] = &entities.Subject{
		ID:          HistoryID,
		Name:        "Historia",
		Description: sp("Análisis de eventos históricos, culturas y sociedades a través del tiempo"),
		Metadata:    sp(`{"area":"Ciencias Sociales","level":"Secundario","difficulty":"Intermedio"}`),
		IsActive:    true,
		CreatedAt:   subjectsBaseTime,
		UpdatedAt:   subjectsBaseTime,
	}

	// Programación
	subjects[ProgrammingID] = &entities.Subject{
		ID:          ProgrammingID,
		Name:        "Programación",
		Description: sp("Introducción a la programación con lenguajes modernos y conceptos de desarrollo de software"),
		Metadata:    sp(`{"area":"Tecnología","level":"Secundario","difficulty":"Avanzado","languages":["Python","Go"]}`),
		IsActive:    true,
		CreatedAt:   subjectsBaseTime,
		UpdatedAt:   subjectsBaseTime,
	}

	// Educación Física
	subjects[PhysicalEdID] = &entities.Subject{
		ID:          PhysicalEdID,
		Name:        "Educación Física",
		Description: sp("Desarrollo de aptitud física, deportes y hábitos saludables"),
		Metadata:    sp(`{"area":"Educación Física","level":"Secundario","difficulty":"Básico"}`),
		IsActive:    true,
		CreatedAt:   subjectsBaseTime,
		UpdatedAt:   subjectsBaseTime,
	}

	return subjects
}

// GetSubjectByName retorna una materia por su nombre
func GetSubjectByName(name string) *entities.Subject {
	subjects := GetSubjects()
	for _, subject := range subjects {
		if subject.Name == name {
			return subject
		}
	}
	return nil
}

// GetActiveSubjects retorna un slice con todas las materias activas
func GetActiveSubjects() []*entities.Subject {
	subjects := GetSubjects()
	var activeSubjects []*entities.Subject

	for _, subject := range subjects {
		if subject.IsActive {
			activeSubjects = append(activeSubjects, subject)
		}
	}

	return activeSubjects
}
