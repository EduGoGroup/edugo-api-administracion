package data

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

var (
	// UUIDs fijos para los materiales de prueba
	MaterialGuidaSumasID  = uuid.MustParse("f1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	MaterialGuiaRestasID  = uuid.MustParse("f2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	MaterialLasPlantasID  = uuid.MustParse("f3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")
	MaterialCicloAguaID   = uuid.MustParse("f4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44")

	// IDs de referencia (school, teachers, academic units)
	schoolPrimariaIDMat   = uuid.MustParse("b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	teacherMathIDMat      = uuid.MustParse("a2eebc99-9c0b-4ef8-bb6d-6bb9bd380a22")
	teacherScienceIDMat   = uuid.MustParse("a3eebc99-9c0b-4ef8-bb6d-6bb9bd380a33")
	seccionPrimer1AMat    = uuid.MustParse("c4eebc99-9c0b-4ef8-bb6d-6bb9bd380a44")
	seccionPrimer1BMat    = uuid.MustParse("c5eebc99-9c0b-4ef8-bb6d-6bb9bd380a55")

	// Timestamp base para consistencia en datos mock
	materialsBaseTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
)

// GetMaterials retorna un mapa con los 4 materiales educativos de prueba
func GetMaterials() map[uuid.UUID]*entities.Material {
	materials := make(map[uuid.UUID]*entities.Material)

	// Helper para crear pointers a strings
	sp := func(s string) *string { return &s }

	// Helper para crear pointers a UUIDs
	up := func(u uuid.UUID) *uuid.UUID { return &u }

	// Guía de Sumas - Matemáticas (PDF)
	materials[MaterialGuidaSumasID] = &entities.Material{
		ID:                    MaterialGuidaSumasID,
		SchoolID:              schoolPrimariaIDMat,
		UploadedByTeacherID:   teacherMathIDMat,
		AcademicUnitID:        up(seccionPrimer1AMat),
		Title:                 "Guía de Sumas",
		Description:           sp("Material educativo sobre sumas básicas"),
		Subject:               sp("Matemáticas"),
		Grade:                 sp("Primaria"),
		FileURL:               "https://s3.example.com/materials/math/suma.pdf",
		FileType:              "application/pdf",
		FileSizeBytes:         1048576,
		Status:                "ready",
		ProcessingStartedAt:   nil,
		ProcessingCompletedAt: nil,
		IsPublic:              true,
		CreatedAt:             materialsBaseTime,
		UpdatedAt:             materialsBaseTime,
		DeletedAt:             nil,
	}

	// Guía de Restas - Matemáticas (PDF)
	materials[MaterialGuiaRestasID] = &entities.Material{
		ID:                    MaterialGuiaRestasID,
		SchoolID:              schoolPrimariaIDMat,
		UploadedByTeacherID:   teacherMathIDMat,
		AcademicUnitID:        up(seccionPrimer1AMat),
		Title:                 "Guía de Restas",
		Description:           sp("Material educativo sobre restas básicas"),
		Subject:               sp("Matemáticas"),
		Grade:                 sp("Primaria"),
		FileURL:               "https://s3.example.com/materials/math/resta.pdf",
		FileType:              "application/pdf",
		FileSizeBytes:         950000,
		Status:                "ready",
		ProcessingStartedAt:   nil,
		ProcessingCompletedAt: nil,
		IsPublic:              true,
		CreatedAt:             materialsBaseTime,
		UpdatedAt:             materialsBaseTime,
		DeletedAt:             nil,
	}

	// Las Plantas - Ciencias (Video MP4)
	materials[MaterialLasPlantasID] = &entities.Material{
		ID:                    MaterialLasPlantasID,
		SchoolID:              schoolPrimariaIDMat,
		UploadedByTeacherID:   teacherScienceIDMat,
		AcademicUnitID:        up(seccionPrimer1BMat),
		Title:                 "Las Plantas",
		Description:           sp("Video educativo sobre plantas"),
		Subject:               sp("Ciencias Naturales"),
		Grade:                 sp("Primaria"),
		FileURL:               "https://s3.example.com/materials/science/plantas.mp4",
		FileType:              "video/mp4",
		FileSizeBytes:         52428800,
		Status:                "ready",
		ProcessingStartedAt:   nil,
		ProcessingCompletedAt: nil,
		IsPublic:              true,
		CreatedAt:             materialsBaseTime,
		UpdatedAt:             materialsBaseTime,
		DeletedAt:             nil,
	}

	// El Ciclo del Agua - Ciencias (PPTX)
	materials[MaterialCicloAguaID] = &entities.Material{
		ID:                    MaterialCicloAguaID,
		SchoolID:              schoolPrimariaIDMat,
		UploadedByTeacherID:   teacherScienceIDMat,
		AcademicUnitID:        up(seccionPrimer1BMat),
		Title:                 "El Ciclo del Agua",
		Description:           sp("Presentación sobre el ciclo del agua"),
		Subject:               sp("Ciencias Naturales"),
		Grade:                 sp("Primaria"),
		FileURL:               "https://s3.example.com/materials/science/agua.pptx",
		FileType:              "application/vnd.ms-powerpoint",
		FileSizeBytes:         2097152,
		Status:                "ready",
		ProcessingStartedAt:   nil,
		ProcessingCompletedAt: nil,
		IsPublic:              true,
		CreatedAt:             materialsBaseTime,
		UpdatedAt:             materialsBaseTime,
		DeletedAt:             nil,
	}

	return materials
}

// GetMaterialByID retorna un material por su ID
func GetMaterialByID(id uuid.UUID) *entities.Material {
	materials := GetMaterials()
	return materials[id]
}

// GetMaterialsByType retorna todos los materiales filtrados por tipo de archivo
func GetMaterialsByType(fileType string) []*entities.Material {
	materials := GetMaterials()
	var filtered []*entities.Material

	for _, material := range materials {
		if material.FileType == fileType {
			filtered = append(filtered, material)
		}
	}

	return filtered
}

// GetMaterialsBySubject retorna todos los materiales filtrados por materia
func GetMaterialsBySubject(subject string) []*entities.Material {
	materials := GetMaterials()
	var filtered []*entities.Material

	for _, material := range materials {
		if material.Subject != nil && *material.Subject == subject {
			filtered = append(filtered, material)
		}
	}

	return filtered
}

// GetMaterialsBySchool retorna todos los materiales de una escuela específica
func GetMaterialsBySchool(schoolID uuid.UUID) []*entities.Material {
	materials := GetMaterials()
	var filtered []*entities.Material

	for _, material := range materials {
		if material.SchoolID == schoolID {
			filtered = append(filtered, material)
		}
	}

	return filtered
}

// GetMaterialsByTeacher retorna todos los materiales subidos por un profesor específico
func GetMaterialsByTeacher(teacherID uuid.UUID) []*entities.Material {
	materials := GetMaterials()
	var filtered []*entities.Material

	for _, material := range materials {
		if material.UploadedByTeacherID == teacherID {
			filtered = append(filtered, material)
		}
	}

	return filtered
}
