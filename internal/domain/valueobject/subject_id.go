package valueobject

import (
	"github.com/edugo/shared/pkg/types"
)

// SubjectID representa el identificador único de una materia
type SubjectID struct {
	value types.UUID
}

// NewSubjectID crea un nuevo SubjectID
func NewSubjectID() SubjectID {
	return SubjectID{value: types.NewUUID()}
}

// SubjectIDFromString crea un SubjectID desde un string
func SubjectIDFromString(s string) (SubjectID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return SubjectID{}, err
	}
	return SubjectID{value: uuid}, nil
}

// String retorna la representación en string
func (s SubjectID) String() string {
	return s.value.String()
}

// UUID retorna el UUID subyacente
func (s SubjectID) UUID() types.UUID {
	return s.value
}

// IsZero verifica si es el valor cero
func (s SubjectID) IsZero() bool {
	return s.value.IsZero()
}
