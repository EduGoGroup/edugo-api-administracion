package valueobject

import (
	"github.com/edugo/shared/pkg/types"
)

// SchoolID representa el identificador único de una escuela
type SchoolID struct {
	value types.UUID
}

// NewSchoolID crea un nuevo SchoolID
func NewSchoolID() SchoolID {
	return SchoolID{value: types.NewUUID()}
}

// SchoolIDFromString crea un SchoolID desde un string
func SchoolIDFromString(s string) (SchoolID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return SchoolID{}, err
	}
	return SchoolID{value: uuid}, nil
}

// String retorna la representación en string
func (s SchoolID) String() string {
	return s.value.String()
}

// UUID retorna el UUID subyacente
func (s SchoolID) UUID() types.UUID {
	return s.value
}

// IsZero verifica si es el valor cero
func (s SchoolID) IsZero() bool {
	return s.value.IsZero()
}

// Equals compara dos SchoolID
func (s SchoolID) Equals(other SchoolID) bool {
	return s.value.String() == other.value.String()
}
