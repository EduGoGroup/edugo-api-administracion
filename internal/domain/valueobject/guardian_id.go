package valueobject

import (
	"github.com/EduGoGroup/edugo-shared/pkg/types"
)

// GuardianID representa el identificador único de un guardian
type GuardianID struct {
	value types.UUID
}

// NewGuardianID crea un nuevo GuardianID
func NewGuardianID() GuardianID {
	return GuardianID{value: types.NewUUID()}
}

// GuardianIDFromString crea un GuardianID desde un string
func GuardianIDFromString(s string) (GuardianID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return GuardianID{}, err
	}
	return GuardianID{value: uuid}, nil
}

// String retorna la representación en string
func (g GuardianID) String() string {
	return g.value.String()
}

// UUID retorna el UUID subyacente
func (g GuardianID) UUID() types.UUID {
	return g.value
}

// IsZero verifica si es el valor cero
func (g GuardianID) IsZero() bool {
	return g.value.IsZero()
}

// Equals compara dos GuardianID
func (g GuardianID) Equals(other GuardianID) bool {
	return g.value.String() == other.value.String()
}
