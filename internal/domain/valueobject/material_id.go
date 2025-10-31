package valueobject

import (
	"github.com/edugo/shared/pkg/types"
)

// MaterialID representa el identificador único de un material
type MaterialID struct {
	value types.UUID
}

// MaterialIDFromString crea un MaterialID desde un string
func MaterialIDFromString(s string) (MaterialID, error) {
	uuid, err := types.ParseUUID(s)
	if err != nil {
		return MaterialID{}, err
	}
	return MaterialID{value: uuid}, nil
}

// String retorna la representación en string
func (m MaterialID) String() string {
	return m.value.String()
}

// IsZero verifica si es el valor cero
func (m MaterialID) IsZero() bool {
	return m.value.IsZero()
}
