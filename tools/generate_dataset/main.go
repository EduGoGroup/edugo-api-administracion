// Package main genera el dataset de mock desde las migraciones SQL de testing.
//
// Uso:
//
//	go run tools/generate_dataset/main.go
//
// Este comando parsea los archivos SQL en edugo-infrastructure/postgres/migrations/testing
// y genera el código Go del dataset en internal/infrastructure/persistence/mock/dataset/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/EduGoGroup/edugo-infrastructure/tools/mock-generator/pkg/generator"
	"github.com/EduGoGroup/edugo-infrastructure/tools/mock-generator/pkg/parser"
)

func main() {
	// Detectar rutas relativas al proyecto
	projectRoot, err := findProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Rutas de entrada y salida
	// Nota: testingDir asume que edugo-infrastructure está clonado como hermano
	// o se puede configurar via variable de entorno
	testingDir := os.Getenv("EDUGO_INFRASTRUCTURE_PATH")
	if testingDir == "" {
		// Asumir estructura de repos-separados
		testingDir = filepath.Join(projectRoot, "..", "edugo-infrastructure", "postgres", "migrations", "testing")
	} else {
		testingDir = filepath.Join(testingDir, "postgres", "migrations", "testing")
	}

	outputDir := filepath.Join(projectRoot, "internal", "infrastructure", "persistence", "mock", "dataset")

	fmt.Println("Mock Dataset Generator")
	fmt.Println("======================")
	fmt.Printf("Testing SQL dir: %s\n", testingDir)
	fmt.Printf("Output dir:      %s\n", outputDir)
	fmt.Println()

	// Verificar que existe el directorio de testing
	if _, err := os.Stat(testingDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: No se encontró el directorio de testing SQL\n")
		fmt.Fprintf(os.Stderr, "Asegúrate de que edugo-infrastructure esté clonado en: %s\n", filepath.Dir(filepath.Dir(testingDir)))
		fmt.Fprintf(os.Stderr, "O configura EDUGO_INFRASTRUCTURE_PATH con la ruta al repo\n")
		os.Exit(1)
	}

	// Parsear archivos SQL
	fmt.Println("Parseando archivos SQL...")
	sp := parser.NewSQLParser()
	tables, err := sp.ParseDirectory(testingDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parseando SQL: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Encontradas %d tablas\n", len(tables))
	for name, table := range tables {
		fmt.Printf("  - %s: %d registros\n", name, len(table.Rows))
	}
	fmt.Println()

	// Generar dataset
	fmt.Println("Generando dataset...")
	gen := generator.NewDatasetGenerator(outputDir, tables)
	if err := gen.Generate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error generando dataset: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Dataset generado exitosamente!")
	fmt.Printf("Archivos en: %s\n", outputDir)
}

// findProjectRoot busca la raíz del proyecto (donde está go.mod)
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no se encontró go.mod en la jerarquía de directorios")
		}
		dir = parent
	}
}
