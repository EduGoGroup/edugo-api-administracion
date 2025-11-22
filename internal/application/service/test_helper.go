package service

import (
	"github.com/EduGoGroup/edugo-shared/logger"
)

// newTestLogger crea un logger para tests que no hace output
func newTestLogger() logger.Logger {
	// Usar logger de consola con nivel error para tests (mínimo output)
	log := logger.NewZapLogger("error", "console")
	return log
}
