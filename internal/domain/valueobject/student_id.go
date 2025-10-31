package valueobject

import (
	"github.com/edugo/shared/pkg/types"
)

// StudentID representa el identificador único de un estudiante
type StudentID struct {
	value types.UUID
}

// NewStudentID crea un nuevo StudentID
func NewStudentID() StudentID {
	return StudentID{value: types.NewUUID()}
}

// StudentIDFromString crea un StudentID desde un string
func StudentIDFromString(s string) (StudentID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return StudentID{}, err
	}
	return StudentID{value: uuid}, nil
}

// String retorna la representación en string
func (s StudentID) String() string {
	return s.value.String()
}

// UUID retorna el UUID subyacente
func (s StudentID) UUID() types.UUID {
	return s.value
}

// IsZero verifica si es el valor cero
func (s StudentID) IsZero() bool {
	return s.value.IsZero()
}

// Equals compara dos StudentID
func (s StudentID) Equals(other StudentID) bool {
	return s.value.String() == other.value.String()
}
