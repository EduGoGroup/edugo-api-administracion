package entity

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAcademicUnit(t *testing.T) {
	schoolID := valueobject.NewSchoolID()
	unitType := valueobject.UnitTypeGrade

	t.Run("should create valid academic unit", func(t *testing.T) {
		unit, err := NewAcademicUnit(schoolID, unitType, "Grade 1", "G1")

		require.NoError(t, err)
		assert.NotNil(t, unit)
		assert.False(t, unit.ID().IsZero())
		assert.Equal(t, schoolID, unit.SchoolID())
		assert.Equal(t, unitType, unit.UnitType())
		assert.Equal(t, "Grade 1", unit.DisplayName())
		assert.Equal(t, "G1", unit.Code())
		assert.True(t, unit.IsRoot())
		assert.False(t, unit.HasChildren())
		assert.NotNil(t, unit.Children())
		assert.Len(t, unit.Children(), 0)
	})

	t.Run("should fail with empty school_id", func(t *testing.T) {
		unit, err := NewAcademicUnit(valueobject.SchoolID{}, unitType, "Grade 1", "G1")

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Contains(t, err.Error(), "school_id is required")
	})

	t.Run("should fail with invalid unit type", func(t *testing.T) {
		invalidType := valueobject.UnitType("invalid")
		unit, err := NewAcademicUnit(schoolID, invalidType, "Grade 1", "G1")

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Contains(t, err.Error(), "invalid unit type")
	})

	t.Run("should fail with empty display name", func(t *testing.T) {
		unit, err := NewAcademicUnit(schoolID, unitType, "", "G1")

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Contains(t, err.Error(), "display_name is required")
	})

	t.Run("should fail with short display name", func(t *testing.T) {
		unit, err := NewAcademicUnit(schoolID, unitType, "AB", "G1")

		assert.Error(t, err)
		assert.Nil(t, unit)
		assert.Contains(t, err.Error(), "display_name must be at least 3 characters")
	})
}

func TestAcademicUnit_SetParent(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should set valid parent", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		err := child.SetParent(parent.ID(), parent.UnitType())

		require.NoError(t, err)
		assert.NotNil(t, child.ParentUnitID())
		assert.True(t, child.ParentUnitID().Equals(parent.ID()))
		assert.False(t, child.IsRoot())
	})

	t.Run("should fail when parent type cannot have children", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "G1-B")

		err := child.SetParent(parent.ID(), parent.UnitType())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot have children")
	})

	t.Run("should fail when unit is its own parent", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.SetParent(unit.ID(), unit.UnitType())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be its own parent")
	})

	t.Run("should fail when child type not allowed", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 2", "G2")

		err := child.SetParent(parent.ID(), parent.UnitType())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be child of")
	})
}

func TestAcademicUnit_RemoveParent(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should remove parent", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = child.SetParent(parent.ID(), parent.UnitType())
		assert.False(t, child.IsRoot())

		child.RemoveParent()
		assert.True(t, child.IsRoot())
		assert.Nil(t, child.ParentUnitID())
	})
}

func TestAcademicUnit_UpdateInfo(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should update display name", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.UpdateInfo("Grade One", "")
		require.NoError(t, err)
		assert.Equal(t, "Grade One", unit.DisplayName())
	})

	t.Run("should update description", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.UpdateInfo("", "New description")
		require.NoError(t, err)
		assert.Equal(t, "New description", unit.Description())
	})

	t.Run("should fail with short display name", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.UpdateInfo("AB", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "display_name must be at least 3 characters")
	})

	t.Run("should fail when both fields are empty", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.UpdateInfo("", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least one field must be provided")
	})
}

func TestAcademicUnit_UpdateDisplayName(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should update display name", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.UpdateDisplayName("Grade One")
		require.NoError(t, err)
		assert.Equal(t, "Grade One", unit.DisplayName())
	})

	t.Run("should fail with empty display name", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.UpdateDisplayName("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "display_name is required")
	})

	t.Run("should fail with short display name", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.UpdateDisplayName("AB")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "display_name must be at least 3 characters")
	})
}

func TestAcademicUnit_CanHaveChildren(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("grade can have children", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		assert.True(t, unit.CanHaveChildren())
	})

	t.Run("section cannot have children", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")
		assert.False(t, unit.CanHaveChildren())
	})

	t.Run("club cannot have children", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeClub, "Chess Club", "CHESS")
		assert.False(t, unit.CanHaveChildren())
	})
}

func TestAcademicUnit_SoftDelete(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should soft delete unit", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		assert.False(t, unit.IsDeleted())

		err := unit.SoftDelete()
		require.NoError(t, err)
		assert.True(t, unit.IsDeleted())
		assert.NotNil(t, unit.DeletedAt())
	})

	t.Run("should fail when already deleted", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		_ = unit.SoftDelete()
		err := unit.SoftDelete()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already deleted")
	})
}

func TestAcademicUnit_Restore(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should restore deleted unit", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		_ = unit.SoftDelete()
		assert.True(t, unit.IsDeleted())

		err := unit.Restore()
		require.NoError(t, err)
		assert.False(t, unit.IsDeleted())
		assert.Nil(t, unit.DeletedAt())
	})

	t.Run("should fail when not deleted", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.Restore()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not deleted")
	})
}

func TestAcademicUnit_Metadata(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should set and get metadata", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		unit.SetMetadata("color", "blue")
		unit.SetMetadata("capacity", 30)

		value, exists := unit.GetMetadata("color")
		assert.True(t, exists)
		assert.Equal(t, "blue", value)

		value, exists = unit.GetMetadata("capacity")
		assert.True(t, exists)
		assert.Equal(t, 30, value)
	})

	t.Run("should return false for non-existent key", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		value, exists := unit.GetMetadata("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, value)
	})
}

// Tree Navigation Tests

func TestAcademicUnit_HasChildren(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should return false when no children", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		assert.False(t, unit.HasChildren())
	})

	t.Run("should return true when has children", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = child.SetParent(parent.ID(), parent.UnitType())
		_ = parent.AddChild(child)

		assert.True(t, parent.HasChildren())
	})
}

func TestAcademicUnit_AddChild(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should add valid child", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = child.SetParent(parent.ID(), parent.UnitType())
		err := parent.AddChild(child)

		require.NoError(t, err)
		assert.True(t, parent.HasChildren())
		assert.Len(t, parent.Children(), 1)
		assert.Equal(t, child.ID(), parent.Children()[0].ID())
	})

	t.Run("should fail when child is nil", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := parent.AddChild(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "child cannot be nil")
	})

	t.Run("should fail when parent cannot have children", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "G1-B")

		err := parent.AddChild(child)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot have children")
	})

	t.Run("should fail when child is same as parent", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.AddChild(unit)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be its own child")
	})

	t.Run("should fail when child has no parent_id", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		err := parent.AddChild(child)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "child must have a parent_id")
	})

	t.Run("should fail when child parent_id does not match", func(t *testing.T) {
		parent1, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		parent2, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 2", "G2")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = child.SetParent(parent2.ID(), parent2.UnitType())
		err := parent1.AddChild(child)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent_id does not match")
	})

	t.Run("should fail when child type not allowed", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 2", "G2")

		_ = child.SetParent(parent.ID(), parent.UnitType())
		err := parent.AddChild(child)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be child of")
	})

	t.Run("should fail when child already added", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = child.SetParent(parent.ID(), parent.UnitType())
		_ = parent.AddChild(child)

		err := parent.AddChild(child)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already added")
	})

	t.Run("should add multiple children", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child1, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")
		child2, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "G1-B")

		_ = child1.SetParent(parent.ID(), parent.UnitType())
		_ = child2.SetParent(parent.ID(), parent.UnitType())

		err1 := parent.AddChild(child1)
		err2 := parent.AddChild(child2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.Len(t, parent.Children(), 2)
	})
}

func TestAcademicUnit_RemoveChild(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should remove child", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = child.SetParent(parent.ID(), parent.UnitType())
		_ = parent.AddChild(child)

		assert.Len(t, parent.Children(), 1)

		err := parent.RemoveChild(child.ID())
		require.NoError(t, err)
		assert.Len(t, parent.Children(), 0)
		assert.False(t, parent.HasChildren())
	})

	t.Run("should fail with zero child_id", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := parent.RemoveChild(valueobject.UnitID{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "child_id is required")
	})

	t.Run("should fail when child not found", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		nonExistentID := valueobject.NewUnitID()

		err := parent.RemoveChild(nonExistentID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "child not found")
	})

	t.Run("should remove specific child from multiple", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child1, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")
		child2, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "G1-B")

		_ = child1.SetParent(parent.ID(), parent.UnitType())
		_ = child2.SetParent(parent.ID(), parent.UnitType())
		_ = parent.AddChild(child1)
		_ = parent.AddChild(child2)

		assert.Len(t, parent.Children(), 2)

		err := parent.RemoveChild(child1.ID())
		require.NoError(t, err)
		assert.Len(t, parent.Children(), 1)
		assert.Equal(t, child2.ID(), parent.Children()[0].ID())
	})
}

func TestAcademicUnit_GetAllDescendants(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should return empty when no children", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		descendants := unit.GetAllDescendants()
		assert.NotNil(t, descendants)
		assert.Len(t, descendants, 0)
	})

	t.Run("should return direct children", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child1, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")
		child2, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "G1-B")

		_ = child1.SetParent(parent.ID(), parent.UnitType())
		_ = child2.SetParent(parent.ID(), parent.UnitType())
		_ = parent.AddChild(child1)
		_ = parent.AddChild(child2)

		descendants := parent.GetAllDescendants()
		assert.Len(t, descendants, 2)
	})

	t.Run("should return all descendants recursively", func(t *testing.T) {
		// Create hierarchy: school -> grade -> section
		school, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSchool, "School", "SCH")
		grade, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		section, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = grade.SetParent(school.ID(), school.UnitType())
		_ = school.AddChild(grade)

		_ = section.SetParent(grade.ID(), grade.UnitType())
		_ = grade.AddChild(section)

		descendants := school.GetAllDescendants()
		assert.Len(t, descendants, 2) // grade + section
		assert.Contains(t, []string{descendants[0].ID().String(), descendants[1].ID().String()}, grade.ID().String())
		assert.Contains(t, []string{descendants[0].ID().String(), descendants[1].ID().String()}, section.ID().String())
	})

	t.Run("should return complex hierarchy", func(t *testing.T) {
		// school -> [grade1, grade2], grade1 -> [sectionA, sectionB]
		school, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSchool, "School", "SCH")
		grade1, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		grade2, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 2", "G2")
		sectionA, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")
		sectionB, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section B", "G1-B")

		_ = grade1.SetParent(school.ID(), school.UnitType())
		_ = grade2.SetParent(school.ID(), school.UnitType())
		_ = school.AddChild(grade1)
		_ = school.AddChild(grade2)

		_ = sectionA.SetParent(grade1.ID(), grade1.UnitType())
		_ = sectionB.SetParent(grade1.ID(), grade1.UnitType())
		_ = grade1.AddChild(sectionA)
		_ = grade1.AddChild(sectionB)

		descendants := school.GetAllDescendants()
		assert.Len(t, descendants, 4) // grade1, grade2, sectionA, sectionB
	})
}

func TestAcademicUnit_GetDepth(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should return 0 for leaf node", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		depth := unit.GetDepth()
		assert.Equal(t, 0, depth)
	})

	t.Run("should return 1 for parent with one level", func(t *testing.T) {
		parent, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		child, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = child.SetParent(parent.ID(), parent.UnitType())
		_ = parent.AddChild(child)

		depth := parent.GetDepth()
		assert.Equal(t, 1, depth)
	})

	t.Run("should return correct depth for multilevel hierarchy", func(t *testing.T) {
		// school -> grade -> section (depth = 2)
		school, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSchool, "School", "SCH")
		grade, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		section, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = grade.SetParent(school.ID(), school.UnitType())
		_ = school.AddChild(grade)

		_ = section.SetParent(grade.ID(), grade.UnitType())
		_ = grade.AddChild(section)

		depth := school.GetDepth()
		assert.Equal(t, 2, depth)
	})

	t.Run("should return max depth when multiple branches", func(t *testing.T) {
		// school -> [grade1 -> section, grade2] (max depth = 2)
		school, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSchool, "School", "SCH")
		grade1, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")
		grade2, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 2", "G2")
		section, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeSection, "Section A", "G1-A")

		_ = grade1.SetParent(school.ID(), school.UnitType())
		_ = grade2.SetParent(school.ID(), school.UnitType())
		_ = school.AddChild(grade1)
		_ = school.AddChild(grade2)

		_ = section.SetParent(grade1.ID(), grade1.UnitType())
		_ = grade1.AddChild(section)

		depth := school.GetDepth()
		assert.Equal(t, 2, depth) // school -> grade1 -> section
	})
}

func TestAcademicUnit_Validate(t *testing.T) {
	schoolID := valueobject.NewSchoolID()

	t.Run("should validate correctly", func(t *testing.T) {
		unit, _ := NewAcademicUnit(schoolID, valueobject.UnitTypeGrade, "Grade 1", "G1")

		err := unit.Validate()
		assert.NoError(t, err)
	})
}

func TestReconstructAcademicUnit(t *testing.T) {
	schoolID := valueobject.NewSchoolID()
	unitID := valueobject.NewUnitID()
	parentID := valueobject.NewUnitID()

	t.Run("should reconstruct unit with all fields", func(t *testing.T) {
		now := time.Now()
		deletedAt := time.Now()
		metadata := map[string]interface{}{"key": "value"}

		unit := ReconstructAcademicUnit(
			unitID,
			&parentID,
			schoolID,
			valueobject.UnitTypeGrade,
			"Grade 1",
			"G1",
			"Description",
			metadata,
			now,
			now,
			&deletedAt,
		)

		assert.NotNil(t, unit)
		assert.Equal(t, unitID, unit.ID())
		assert.NotNil(t, unit.ParentUnitID())
		assert.Equal(t, parentID, *unit.ParentUnitID())
		assert.Equal(t, schoolID, unit.SchoolID())
		assert.Equal(t, valueobject.UnitTypeGrade, unit.UnitType())
		assert.Equal(t, "Grade 1", unit.DisplayName())
		assert.Equal(t, "G1", unit.Code())
		assert.Equal(t, "Description", unit.Description())
		assert.True(t, unit.IsDeleted())
		assert.NotNil(t, unit.Children())
		assert.Len(t, unit.Children(), 0)
	})

	t.Run("should reconstruct unit with nil metadata", func(t *testing.T) {
		now := time.Now()

		unit := ReconstructAcademicUnit(
			unitID,
			nil,
			schoolID,
			valueobject.UnitTypeGrade,
			"Grade 1",
			"G1",
			"",
			nil,
			now,
			now,
			nil,
		)

		assert.NotNil(t, unit)
		assert.NotNil(t, unit.Metadata())
		assert.Len(t, unit.Metadata(), 0)
	})
}
