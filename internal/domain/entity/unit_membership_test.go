package entity

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUnitMembership(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()
	validFrom := time.Now()

	t.Run("should create valid membership", func(t *testing.T) {
		membership, err := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		require.NoError(t, err)
		assert.NotNil(t, membership)
		assert.False(t, membership.ID().IsZero())
		assert.Equal(t, unitID, membership.UnitID())
		assert.Equal(t, userID, membership.UserID())
		assert.Equal(t, valueobject.RoleStudent, membership.Role())
		assert.Equal(t, validFrom.Unix(), membership.ValidFrom().Unix())
		assert.Nil(t, membership.ValidUntil())
		assert.True(t, membership.IsActive())
	})

	t.Run("should fail with empty unit_id", func(t *testing.T) {
		membership, err := NewUnitMembership(valueobject.UnitID{}, userID, valueobject.RoleStudent, validFrom)

		assert.Error(t, err)
		assert.Nil(t, membership)
		assert.Contains(t, err.Error(), "unit_id is required")
	})

	t.Run("should fail with empty user_id", func(t *testing.T) {
		membership, err := NewUnitMembership(unitID, valueobject.UserID{}, valueobject.RoleStudent, validFrom)

		assert.Error(t, err)
		assert.Nil(t, membership)
		assert.Contains(t, err.Error(), "user_id is required")
	})

	t.Run("should fail with invalid role", func(t *testing.T) {
		invalidRole := valueobject.MembershipRole("invalid")
		membership, err := NewUnitMembership(unitID, userID, invalidRole, validFrom)

		assert.Error(t, err)
		assert.Nil(t, membership)
		assert.Contains(t, err.Error(), "invalid membership role")
	})

	t.Run("should use current time if validFrom is zero", func(t *testing.T) {
		before := time.Now()
		membership, err := NewUnitMembership(unitID, userID, valueobject.RoleStudent, time.Time{})
		after := time.Now()

		require.NoError(t, err)
		assert.NotNil(t, membership)
		assert.False(t, membership.ValidFrom().Before(before))
		assert.False(t, membership.ValidFrom().After(after))
	})
}

func TestUnitMembership_IsActive(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should be active when no validUntil", func(t *testing.T) {
		validFrom := time.Now().Add(-24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		assert.True(t, membership.IsActive())
	})

	t.Run("should be active when current time is before validUntil", func(t *testing.T) {
		validFrom := time.Now().Add(-24 * time.Hour)
		validUntil := time.Now().Add(24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		_ = membership.SetValidUntil(validUntil)

		assert.True(t, membership.IsActive())
	})

	t.Run("should not be active when current time is after validUntil", func(t *testing.T) {
		validFrom := time.Now().Add(-48 * time.Hour)
		validUntil := time.Now().Add(-24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		_ = membership.SetValidUntil(validUntil)

		assert.False(t, membership.IsActive())
	})

	t.Run("should not be active when current time is before validFrom", func(t *testing.T) {
		validFrom := time.Now().Add(24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		assert.False(t, membership.IsActive())
	})
}

func TestUnitMembership_IsActiveAt(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should be active at specific time within range", func(t *testing.T) {
		validFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		validUntil := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		checkTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)

		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		_ = membership.SetValidUntil(validUntil)

		assert.True(t, membership.IsActiveAt(checkTime))
	})

	t.Run("should not be active before validFrom", func(t *testing.T) {
		validFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		checkTime := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		assert.False(t, membership.IsActiveAt(checkTime))
	})

	t.Run("should not be active after validUntil", func(t *testing.T) {
		validFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		validUntil := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		checkTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		_ = membership.SetValidUntil(validUntil)

		assert.False(t, membership.IsActiveAt(checkTime))
	})

	t.Run("should be active at validUntil exact time", func(t *testing.T) {
		validFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		validUntil := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		_ = membership.SetValidUntil(validUntil)

		assert.True(t, membership.IsActiveAt(validUntil))
	})
}

func TestUnitMembership_SetValidUntil(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should set valid until", func(t *testing.T) {
		validFrom := time.Now().Add(-24 * time.Hour)
		validUntil := time.Now().Add(24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		err := membership.SetValidUntil(validUntil)

		require.NoError(t, err)
		assert.NotNil(t, membership.ValidUntil())
		assert.Equal(t, validUntil.Unix(), membership.ValidUntil().Unix())
	})

	t.Run("should fail when validUntil is before validFrom", func(t *testing.T) {
		validFrom := time.Now()
		validUntil := time.Now().Add(-24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		err := membership.SetValidUntil(validUntil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "valid_until must be after valid_from")
	})

	t.Run("should fail when validUntil equals validFrom", func(t *testing.T) {
		validFrom := time.Now()
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		err := membership.SetValidUntil(validFrom)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "valid_until must be after valid_from")
	})
}

func TestUnitMembership_ExtendIndefinitely(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should extend indefinitely", func(t *testing.T) {
		validFrom := time.Now().Add(-24 * time.Hour)
		validUntil := time.Now().Add(24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		_ = membership.SetValidUntil(validUntil)

		assert.NotNil(t, membership.ValidUntil())

		membership.ExtendIndefinitely()

		assert.Nil(t, membership.ValidUntil())
		assert.True(t, membership.IsActive())
	})
}

func TestUnitMembership_Expire(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should expire membership", func(t *testing.T) {
		validFrom := time.Now().Add(-24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		assert.True(t, membership.IsActive())

		err := membership.Expire()

		require.NoError(t, err)
		assert.NotNil(t, membership.ValidUntil())
		assert.False(t, membership.IsActive())
	})

	t.Run("should fail when already expired", func(t *testing.T) {
		validFrom := time.Now().Add(-48 * time.Hour)
		validUntil := time.Now().Add(-24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		_ = membership.SetValidUntil(validUntil)

		err := membership.Expire()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already expired")
	})
}

func TestUnitMembership_ChangeRole(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should change role", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, time.Now())

		assert.Equal(t, valueobject.RoleStudent, membership.Role())

		err := membership.ChangeRole(valueobject.RoleTeacher)

		require.NoError(t, err)
		assert.Equal(t, valueobject.RoleTeacher, membership.Role())
	})

	t.Run("should fail with invalid role", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, time.Now())

		invalidRole := valueobject.MembershipRole("invalid")
		err := membership.ChangeRole(invalidRole)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid membership role")
	})

	t.Run("should fail when role is same", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, time.Now())

		err := membership.ChangeRole(valueobject.RoleStudent)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role is already")
	})
}

func TestUnitMembership_HasPermission(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("admin should have all permissions", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleAdmin, time.Now().Add(-1*time.Hour))

		assert.True(t, membership.HasPermission("manage_unit"))
		assert.True(t, membership.HasPermission("view_members"))
		assert.True(t, membership.HasPermission("add_members"))
		assert.True(t, membership.HasPermission("any_permission"))
	})

	t.Run("coordinator should have management permissions", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleCoordinator, time.Now().Add(-1*time.Hour))

		assert.True(t, membership.HasPermission("manage_unit"))
		assert.True(t, membership.HasPermission("view_members"))
		assert.True(t, membership.HasPermission("add_members"))
		assert.False(t, membership.HasPermission("other_permission"))
	})

	t.Run("teacher should have view permissions", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleTeacher, time.Now().Add(-1*time.Hour))

		assert.True(t, membership.HasPermission("view_members"))
		assert.False(t, membership.HasPermission("manage_unit"))
		assert.False(t, membership.HasPermission("add_members"))
	})

	t.Run("student should have no special permissions", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, time.Now().Add(-1*time.Hour))

		assert.False(t, membership.HasPermission("view_members"))
		assert.False(t, membership.HasPermission("manage_unit"))
		assert.False(t, membership.HasPermission("add_members"))
	})

	t.Run("inactive membership should have no permissions", func(t *testing.T) {
		validFrom := time.Now().Add(-48 * time.Hour)
		validUntil := time.Now().Add(-24 * time.Hour)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleAdmin, validFrom)
		_ = membership.SetValidUntil(validUntil)

		assert.False(t, membership.IsActive())
		assert.False(t, membership.HasPermission("manage_unit"))
		assert.False(t, membership.HasPermission("any_permission"))
	})
}

func TestUnitMembership_Metadata(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should set and get metadata", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, time.Now())

		membership.SetMetadata("notes", "Good student")
		membership.SetMetadata("attendance", 95.5)

		value, exists := membership.GetMetadata("notes")
		assert.True(t, exists)
		assert.Equal(t, "Good student", value)

		value, exists = membership.GetMetadata("attendance")
		assert.True(t, exists)
		assert.Equal(t, 95.5, value)
	})

	t.Run("should return false for non-existent key", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, time.Now())

		value, exists := membership.GetMetadata("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, value)
	})
}

func TestUnitMembership_Validate(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should validate correctly", func(t *testing.T) {
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, time.Now())

		err := membership.Validate()
		assert.NoError(t, err)
	})

	t.Run("should fail validation when validUntil before validFrom", func(t *testing.T) {
		validFrom := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		// Manually set invalid validUntil (bypass SetValidUntil validation for testing)
		invalidUntil := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		membership = ReconstructUnitMembership(
			membership.ID(),
			unitID,
			userID,
			valueobject.RoleStudent,
			validFrom,
			&invalidUntil,
			nil,
			membership.CreatedAt(),
			membership.UpdatedAt(),
		)

		err := membership.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "valid_until must be after valid_from")
	})
}

func TestReconstructUnitMembership(t *testing.T) {
	membershipID := valueobject.NewMembershipID()
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should reconstruct membership with all fields", func(t *testing.T) {
		validFrom := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		validUntil := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
		createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		metadata := map[string]interface{}{"key": "value"}

		membership := ReconstructUnitMembership(
			membershipID,
			unitID,
			userID,
			valueobject.RoleStudent,
			validFrom,
			&validUntil,
			metadata,
			createdAt,
			updatedAt,
		)

		assert.NotNil(t, membership)
		assert.Equal(t, membershipID, membership.ID())
		assert.Equal(t, unitID, membership.UnitID())
		assert.Equal(t, userID, membership.UserID())
		assert.Equal(t, valueobject.RoleStudent, membership.Role())
		assert.Equal(t, validFrom, membership.ValidFrom())
		assert.NotNil(t, membership.ValidUntil())
		assert.Equal(t, validUntil, *membership.ValidUntil())
	})

	t.Run("should reconstruct membership with nil metadata", func(t *testing.T) {
		validFrom := time.Now()
		createdAt := time.Now()
		updatedAt := time.Now()

		membership := ReconstructUnitMembership(
			membershipID,
			unitID,
			userID,
			valueobject.RoleStudent,
			validFrom,
			nil,
			nil,
			createdAt,
			updatedAt,
		)

		assert.NotNil(t, membership)
		assert.NotNil(t, membership.Metadata())
		assert.Len(t, membership.Metadata(), 0)
	})
}

func TestUnitMembership_Getters(t *testing.T) {
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()
	validFrom := time.Now()

	membership, _ := NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

	t.Run("should return correct ID", func(t *testing.T) {
		assert.False(t, membership.ID().IsZero())
	})

	t.Run("should return correct UnitID", func(t *testing.T) {
		assert.Equal(t, unitID, membership.UnitID())
	})

	t.Run("should return correct UserID", func(t *testing.T) {
		assert.Equal(t, userID, membership.UserID())
	})

	t.Run("should return correct Role", func(t *testing.T) {
		assert.Equal(t, valueobject.RoleStudent, membership.Role())
	})

	t.Run("should return correct ValidFrom", func(t *testing.T) {
		assert.Equal(t, validFrom.Unix(), membership.ValidFrom().Unix())
	})

	t.Run("should return correct ValidUntil", func(t *testing.T) {
		assert.Nil(t, membership.ValidUntil())
	})

	t.Run("should return correct timestamps", func(t *testing.T) {
		assert.False(t, membership.CreatedAt().IsZero())
		assert.False(t, membership.UpdatedAt().IsZero())
	})

	t.Run("should return metadata copy", func(t *testing.T) {
		membership.SetMetadata("test", "value")
		metadata := membership.Metadata()

		// Modify the copy
		metadata["test"] = "modified"

		// Original should not be affected
		value, _ := membership.GetMetadata("test")
		assert.Equal(t, "value", value)
	})
}
