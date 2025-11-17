//go:build tools
// +build tools

package integration

// Este archivo asegura que go.mod mantenga la dependencia de infrastructure
// aunque no importemos el paquete directamente (porque es un main package)
import (
	_ "github.com/EduGoGroup/edugo-infrastructure/postgres"
)
