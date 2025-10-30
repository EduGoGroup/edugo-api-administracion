package valueobject

import (
	"github.com/EduGoGroup/edugo-shared/pkg/types"
)

// UnitID representa el identificador único de una unidad
type UnitID struct {
	value types.UUID
}

// NewUnitID crea un nuevo UnitID
func NewUnitID() UnitID {
	return UnitID{value: types.NewUUID()}
}

// UnitIDFromString crea un UnitID desde un string
func UnitIDFromString(s string) (UnitID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return UnitID{}, err
	}
	return UnitID{value: uuid}, nil
}

// String retorna la representación en string
func (u UnitID) String() string {
	return u.value.String()
}

// UUID retorna el UUID subyacente
func (u UnitID) UUID() types.UUID {
	return u.value
}

// IsZero verifica si es el valor cero
func (u UnitID) IsZero() bool {
	return u.value.IsZero()
}

// Equals compara dos UnitID
func (u UnitID) Equals(other UnitID) bool {
	return u.value.String() == other.value.String()
}
