package valueobject

import (
	"github.com/EduGoGroup/edugo-shared/common/types"
)

// MembershipID representa el identificador único de una membresía
type MembershipID struct {
	value types.UUID
}

// NewMembershipID crea un nuevo MembershipID
func NewMembershipID() MembershipID {
	return MembershipID{value: types.NewUUID()}
}

// MembershipIDFromString crea un MembershipID desde un string
func MembershipIDFromString(s string) (MembershipID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return MembershipID{}, err
	}
	return MembershipID{value: uuid}, nil
}

// String retorna la representación en string
func (m MembershipID) String() string {
	return m.value.String()
}

// UUID retorna el UUID subyacente
func (m MembershipID) UUID() types.UUID {
	return m.value
}

// IsZero verifica si es el valor cero
func (m MembershipID) IsZero() bool {
	return m.value.IsZero()
}

// Equals compara dos MembershipID
func (m MembershipID) Equals(other MembershipID) bool {
	return m.value.String() == other.value.String()
}
