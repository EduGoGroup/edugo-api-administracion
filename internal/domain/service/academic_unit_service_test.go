package service

import (
	"testing"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAcademicUnitService_SetParent(t *testing.T) {
	service := NewAcademicUnitDomainService()
	schoolID := valueobject.NewSchoolID()

	t.Run("should set valid parent", func(t *testing.T) {
		parent, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		err := service.SetParent(child, parent.ID(), parent.UnitType())

		assert.NoError(t, err)
		assert.NotNil(t, child.ParentUnitID())
		assert.True(t, child.ParentUnitID().Equals(parent.ID()))
	})

	t.Run("should fail when parent type cannot have children", func(t *testing.T) {
		parent, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "S-A")
		child, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "S-B")

		err := service.SetParent(child, parent.ID(), parent.UnitType())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot have children")
	})

	t.Run("should fail when unit is its own parent", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := service.SetParent(unit, unit.ID(), unit.UnitType())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be its own parent")
	})

	t.Run("should fail when child type not allowed", func(t *testing.T) {
		parent, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 2", "G2")

		err := service.SetParent(child, parent.ID(), parent.UnitType())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be child of")
	})
}

func TestAcademicUnitService_AddChild(t *testing.T) {
	service := NewAcademicUnitDomainService()
	schoolID := valueobject.NewSchoolID()

	t.Run("should add valid child", func(t *testing.T) {
		parent, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		// Primero establecer la relación padre-hijo
		err := service.SetParent(child, parent.ID(), parent.UnitType())
		require.NoError(t, err)

		// Luego agregar el hijo al padre
		err = service.AddChild(parent, child)
		assert.NoError(t, err)
		assert.True(t, service.HasChildren(parent))
	})

	t.Run("should fail when child is nil", func(t *testing.T) {
		parent, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := service.AddChild(parent, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "child cannot be nil")
	})

	t.Run("should fail when parent cannot have children", func(t *testing.T) {
		parent, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "S-A")
		child, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "S-B")

		err := service.AddChild(parent, child)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot have children")
	})
}

func TestAcademicUnitService_GetAllDescendants(t *testing.T) {
	service := NewAcademicUnitDomainService()
	schoolID := valueobject.NewSchoolID()

	t.Run("should return empty when no children", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		descendants := service.GetAllDescendants(unit)

		assert.Empty(t, descendants)
	})

	t.Run("should return all descendants recursively", func(t *testing.T) {
		// Crear jerarquía: Grade -> Section1, Section2 -> (Section1 tiene subsección si fuera posible)
		grade, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		section1, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")
		section2, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "G1-B")

		_ = service.SetParent(section1, grade.ID(), grade.UnitType())
		_ = service.SetParent(section2, grade.ID(), grade.UnitType())
		_ = service.AddChild(grade, section1)
		_ = service.AddChild(grade, section2)

		descendants := service.GetAllDescendants(grade)

		assert.Len(t, descendants, 2)
	})
}

func TestAcademicUnitService_GetDepth(t *testing.T) {
	service := NewAcademicUnitDomainService()
	schoolID := valueobject.NewSchoolID()

	t.Run("should return 0 for leaf node", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "S-A")

		depth := service.GetDepth(unit)

		assert.Equal(t, 0, depth)
	})

	t.Run("should return correct depth for hierarchy", func(t *testing.T) {
		grade, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		section, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = service.SetParent(section, grade.ID(), grade.UnitType())
		_ = service.AddChild(grade, section)

		depth := service.GetDepth(grade)

		assert.Equal(t, 1, depth)
	})
}

func TestAcademicUnitService_UpdateDisplayName(t *testing.T) {
	service := NewAcademicUnitDomainService()
	schoolID := valueobject.NewSchoolID()

	t.Run("should update display name", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := service.UpdateDisplayName(unit, "First Grade")

		assert.NoError(t, err)
		assert.Equal(t, "First Grade", unit.DisplayName())
	})

	t.Run("should fail with empty display name", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := service.UpdateDisplayName(unit, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "display_name is required")
	})

	t.Run("should fail with short display name", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := service.UpdateDisplayName(unit, "AB")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be at least 3 characters")
	})
}

func TestAcademicUnitService_SoftDelete(t *testing.T) {
	service := NewAcademicUnitDomainService()
	schoolID := valueobject.NewSchoolID()

	t.Run("should soft delete unit", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := service.SoftDelete(unit)

		assert.NoError(t, err)
		assert.True(t, unit.IsDeleted())
	})

	t.Run("should fail when already deleted", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		_ = service.SoftDelete(unit)

		err := service.SoftDelete(unit)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already deleted")
	})
}

func TestAcademicUnitService_Restore(t *testing.T) {
	service := NewAcademicUnitDomainService()
	schoolID := valueobject.NewSchoolID()

	t.Run("should restore deleted unit", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		_ = service.SoftDelete(unit)

		err := service.Restore(unit)

		assert.NoError(t, err)
		assert.False(t, unit.IsDeleted())
	})

	t.Run("should fail when not deleted", func(t *testing.T) {
		unit, _ := entity.NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := service.Restore(unit)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not deleted")
	})
}
