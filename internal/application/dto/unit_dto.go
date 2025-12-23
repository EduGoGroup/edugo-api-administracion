package dto

import (
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
)

type CreateUnitRequest struct {
	SchoolID     string  `json:"school_id" binding:"required,uuid" validate:"required,uuid"`
	ParentUnitID *string `json:"parent_unit_id" binding:"omitempty,uuid" validate:"omitempty,uuid"`
	Name         string  `json:"name" binding:"required,min=2" validate:"required,min=2"`
	Description  string  `json:"description"`
}

type UpdateUnitRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=2" validate:"omitempty,min=2"`
	Description *string `json:"description"`
}

type UnitResponse struct {
	ID           string    `json:"id"`
	SchoolID     string    `json:"school_id"`
	ParentUnitID *string   `json:"parent_unit_id,omitempty"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func ToUnitResponse(unit *entities.Unit) UnitResponse {
	var parentID *string
	if unit.ParentUnitID != nil {
		pid := unit.ParentUnitID.String()
		parentID = &pid
	}
	desc := ""
	if unit.Description != nil {
		desc = *unit.Description
	}
	return UnitResponse{
		ID:           unit.ID.String(),
		SchoolID:     unit.SchoolID.String(),
		ParentUnitID: parentID,
		Name:         unit.Name,
		Description:  desc,
		IsActive:     unit.IsActive,
		CreatedAt:    unit.CreatedAt,
		UpdatedAt:    unit.UpdatedAt,
	}
}
