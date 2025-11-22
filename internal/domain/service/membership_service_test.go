package service

import (
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/entity"
	"github.com/EduGoGroup/edugo-api-administracion/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
)

func TestMembershipService_IsActive(t *testing.T) {
	service := NewMembershipDomainService()
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()
	validFrom := time.Now().Add(-24 * time.Hour)

	t.Run("should be active when no validUntil", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		isActive := service.IsActive(membership)

		assert.True(t, isActive)
	})

	t.Run("should be active when current time is before validUntil", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		validUntil := time.Now().Add(24 * time.Hour)
		_ = service.SetValidUntil(membership, validUntil)

		isActive := service.IsActive(membership)

		assert.True(t, isActive)
	})

	t.Run("should not be active when current time is after validUntil", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		validUntil := time.Now().Add(-1 * time.Hour)
		membership.SetValidUntilValue(&validUntil)

		isActive := service.IsActive(membership)

		assert.False(t, isActive)
	})

	t.Run("should not be active when current time is before validFrom", func(t *testing.T) {
		futureValidFrom := time.Now().Add(24 * time.Hour)
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, futureValidFrom)

		isActive := service.IsActive(membership)

		assert.False(t, isActive)
	})
}

func TestMembershipService_IsActiveAt(t *testing.T) {
	service := NewMembershipDomainService()
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()

	t.Run("should be active at specific time within range", func(t *testing.T) {
		validFrom := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		validUntil := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
		checkTime := time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC)

		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		membership.SetValidUntilValue(&validUntil)

		isActive := service.IsActiveAt(membership, checkTime)

		assert.True(t, isActive)
	})

	t.Run("should not be active before validFrom", func(t *testing.T) {
		validFrom := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		checkTime := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		isActive := service.IsActiveAt(membership, checkTime)

		assert.False(t, isActive)
	})

	t.Run("should not be active after validUntil", func(t *testing.T) {
		validFrom := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		validUntil := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
		checkTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		membership.SetValidUntilValue(&validUntil)

		isActive := service.IsActiveAt(membership, checkTime)

		assert.False(t, isActive)
	})
}

func TestMembershipService_SetValidUntil(t *testing.T) {
	service := NewMembershipDomainService()
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()
	validFrom := time.Now()

	t.Run("should set valid until", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		validUntil := validFrom.Add(365 * 24 * time.Hour)

		err := service.SetValidUntil(membership, validUntil)

		assert.NoError(t, err)
		assert.NotNil(t, membership.ValidUntil())
		assert.Equal(t, validUntil.Unix(), membership.ValidUntil().Unix())
	})

	t.Run("should fail when validUntil is before validFrom", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		validUntil := validFrom.Add(-24 * time.Hour)

		err := service.SetValidUntil(membership, validUntil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be after")
	})

	t.Run("should fail when validUntil equals validFrom", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		err := service.SetValidUntil(membership, validFrom)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be equal")
	})
}

func TestMembershipService_Expire(t *testing.T) {
	service := NewMembershipDomainService()
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()
	validFrom := time.Now().Add(-24 * time.Hour)

	t.Run("should expire membership", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		err := service.Expire(membership)

		assert.NoError(t, err)
		assert.NotNil(t, membership.ValidUntil())
		assert.False(t, service.IsActive(membership))
	})

	t.Run("should fail when already expired", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		_ = service.Expire(membership)

		err := service.Expire(membership)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already expired")
	})
}

func TestMembershipService_ChangeRole(t *testing.T) {
	service := NewMembershipDomainService()
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()
	validFrom := time.Now()

	t.Run("should change role", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		err := service.ChangeRole(membership, valueobject.RoleTeacher)

		assert.NoError(t, err)
		assert.Equal(t, valueobject.RoleTeacher, membership.Role())
	})

	t.Run("should fail with invalid role", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)
		invalidRole := valueobject.MembershipRole("invalid")

		err := service.ChangeRole(membership, invalidRole)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
	})

	t.Run("should fail when role is same", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleStudent, validFrom)

		err := service.ChangeRole(membership, valueobject.RoleStudent)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be different")
	})
}

func TestMembershipService_HasPermission(t *testing.T) {
	service := NewMembershipDomainService()
	unitID := valueobject.NewUnitID()
	userID := valueobject.NewUserID()
	validFrom := time.Now()

	t.Run("admin should have all permissions", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleAdmin, validFrom)

		assert.True(t, service.HasPermission(membership, "view"))
		assert.True(t, service.HasPermission(membership, "edit"))
		assert.True(t, service.HasPermission(membership, "delete"))
		assert.True(t, service.HasPermission(membership, "manage_members"))
	})

	t.Run("coordinator should have management permissions", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleCoordinator, validFrom)

		assert.True(t, service.HasPermission(membership, "view"))
		assert.True(t, service.HasPermission(membership, "edit"))
		assert.True(t, service.HasPermission(membership, "manage_members"))
		assert.False(t, service.HasPermission(membership, "delete"))
	})

	t.Run("teacher should have view permissions", func(t *testing.T) {
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleTeacher, validFrom)

		assert.True(t, service.HasPermission(membership, "view"))
		assert.True(t, service.HasPermission(membership, "edit"))
		assert.False(t, service.HasPermission(membership, "delete"))
		assert.False(t, service.HasPermission(membership, "manage_members"))
	})

	t.Run("inactive membership should have no permissions", func(t *testing.T) {
		pastValidFrom := time.Now().Add(-48 * time.Hour)
		membership, _ := entity.NewUnitMembership(unitID, userID, valueobject.RoleAdmin, pastValidFrom)
		_ = service.Expire(membership)

		assert.False(t, service.HasPermission(membership, "view"))
		assert.False(t, service.HasPermission(membership, "edit"))
	})
}
