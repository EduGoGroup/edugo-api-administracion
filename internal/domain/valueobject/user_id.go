package valueobject

import (
	"github.com/EduGoGroup/edugo-shared/pkg/types"
)

// UserID representa el identificador único de un usuario
type UserID struct {
	value types.UUID
}

// NewUserID crea un nuevo UserID
func NewUserID() UserID {
	return UserID{value: types.NewUUID()}
}

// UserIDFromString crea un UserID desde un string
func UserIDFromString(s string) (UserID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return UserID{}, err
	}
	return UserID{value: uuid}, nil
}

// String retorna la representación en string
func (u UserID) String() string {
	return u.value.String()
}

// UUID retorna el UUID subyacente
func (u UserID) UUID() types.UUID {
	return u.value
}

// IsZero verifica si es el valor cero
func (u UserID) IsZero() bool {
	return u.value.IsZero()
}

// Equals compara dos UserID
func (u UserID) Equals(other UserID) bool {
	return u.value.String() == other.value.String()
}
