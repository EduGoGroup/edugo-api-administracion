# Mock Repositories

Implementación mock (in-memory) de repositorios para `edugo-api-administracion`.

## 🎯 ¿Por qué Mock Repositories?

Los repositorios mock permiten ejecutar la API **sin infraestructura externa** (PostgreSQL, MongoDB, RabbitMQ):

### Beneficios

- **🚀 Desarrollo Frontend**: No necesita Docker/PostgreSQL corriendo
- **💾 Ahorro de RAM**: ≈1.2GB menos (sin PostgreSQL ~500MB, MongoDB ~400MB, RabbitMQ ~300MB)
- **⚡ Startup Rápido**: API inicia en <3 segundos
- **🧪 Testing**: Datos predecibles y reseteables
- **📦 Portabilidad**: API funciona sin configurar infraestructura
- **🔄 Consistencia**: Mismo dataset siempre disponible

## 📋 Uso

### Activar Mocks

**Opción 1: Archivo de Configuración**

```yaml
# config/config-local.yaml
database:
  use_mock_repositories: true
```

**Opción 2: Variable de Entorno**

```bash
export USE_MOCK_REPOSITORIES=true
make run
```

### Desactivar Mocks (usar PostgreSQL)

```yaml
# config/config.yaml
database:
  use_mock_repositories: false
```

## 📊 Datos Demo Disponibles (Sprints 1, 2 y 3 - COMPLETO)

### Users (8 usuarios)

| Email | Rol | Nombre | Contraseña |
|-------|-----|--------|------------|
| `admin@edugo.test` | admin | Admin Demo | `edugo2024` |
| `teacher.math@edugo.test` | teacher | María García | `edugo2024` |
| `teacher.science@edugo.test` | teacher | Juan Pérez | `edugo2024` |
| `student1@edugo.test` | student | Carlos Rodríguez | `edugo2024` |
| `student2@edugo.test` | student | Ana Martínez | `edugo2024` |
| `student3@edugo.test` | student | Luis González | `edugo2024` |
| `guardian1@edugo.test` | guardian | Roberto Fernández | `edugo2024` |
| `guardian2@edugo.test` | guardian | Patricia López | `edugo2024` |

**Nota**: Todos los usuarios tienen la misma contraseña para facilitar el testing.

### Schools (3 escuelas)

| Código | Nombre | Ciudad | Tier |
|--------|--------|--------|------|
| `SCH_PRI_001` | Escuela Primaria Demo | Buenos Aires | basic |
| `SCH_SEC_001` | Colegio Secundario Demo | Buenos Aires | premium |
| `SCH_TEC_001` | Instituto Técnico Demo | Córdoba | premium |

### Academic Units (12 unidades jerárquicas)

Estructura jerárquica completa con grados y secciones:
- **Escuela Primaria**: Primer Grado (Sección A, B), Segundo Grado, Tercer Grado
- **Colegio Secundario**: Primer Año (Sección 1, 2), Segundo Año
- **Instituto Técnico**: Programación I, Bases de Datos

### Memberships (5 asignaciones)

- Teacher María → Escuela Primaria → Sección A
- Teacher Juan → Escuela Primaria → Sección B  
- Student Carlos → Escuela Primaria → Sección A
- Student Ana → Escuela Primaria → Sección A
- Student Luis → Escuela Primaria → Sección B

### Subjects (6 materias)

- Matemáticas, Ciencias Naturales, Lengua y Literatura
- Historia, Programación, Educación Física

### Units (4 unidades organizacionales)

- Departamento de Matemáticas, Departamento de Ciencias
- Coordinación Académica, Grupo de Docentes

### Materials (4 materiales educativos)

- Guía de Sumas (PDF), Guía de Restas (PDF)
- Las Plantas (Video MP4), El Ciclo del Agua (PPTX)

### Guardian Relations (3 relaciones)

- Roberto → Carlos (father)
- Patricia → Ana (mother)
- Roberto → Luis (legal_guardian)

## 🏗️ Arquitectura

```
mock/
├── README.md              # Este archivo
├── data/                  # Datos estáticos pre-cargados (8 archivos)
│   ├── users.go          # 8 usuarios demo
│   ├── schools.go        # 3 escuelas demo
│   ├── academic_units.go # 12 unidades jerárquicas
│   ├── memberships.go    # 5 memberships
│   ├── subjects.go       # 6 materias
│   ├── units.go          # 4 unidades organizacionales
│   ├── materials.go      # 4 materiales educativos
│   └── guardian_relations.go # 3 relaciones tutor-estudiante
└── repository/            # Implementaciones mock (9 repositorios)
    ├── school_repository_mock.go
    ├── user_repository_mock.go
    ├── academic_unit_repository_mock.go
    ├── unit_membership_repository_mock.go
    ├── subject_repository_mock.go
    ├── unit_repository_mock.go
    ├── material_repository_mock.go
    ├── guardian_repository_mock.go
    └── stats_repository_mock.go
```

### Características Técnicas

✅ **Thread-safe**: Todos los mocks usan `sync.RWMutex`  
✅ **Inmutables**: Retornan copias, no referencias  
✅ **Validaciones**: Replican comportamiento PostgreSQL  
✅ **Soft Delete**: Respetan campo `DeletedAt`  
✅ **Errores consistentes**: Mismos tipos que implementación real  

## 🧪 Testing

### Ejemplo de Uso en Tests

```go
package test

import (
    "context"
    "testing"
    
    "github.com/EduGoGroup/edugo-api-administracion/internal/factory"
    "github.com/stretchr/testify/assert"
)

func TestSchoolCRUD(t *testing.T) {
    // Crear factory mock
    factory := factory.NewMockRepositoryFactory()
    repo := factory.CreateSchoolRepository()
    
    ctx := context.Background()
    
    // Buscar escuela demo
    school, err := repo.FindByCode(ctx, "SCH_PRI_001")
    assert.NoError(t, err)
    assert.Equal(t, "Escuela Primaria Demo", school.Name)
}
```

### Login de Prueba

```bash
# Iniciar API con mocks
make run

# Login con usuario admin
curl -X POST http://localhost:8081/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@edugo.test",
    "password": "edugo2024"
  }'
```

## 🚦 Estado de Implementación

### Sprint 1 (MVP Core) - ✅ COMPLETADO

- [x] Factory Pattern (RepositoryFactory)
- [x] PostgresFactory
- [x] MockFactory
- [x] SchoolRepository Mock (9 métodos)
- [x] UserRepository Mock (7 métodos)
- [x] Datos mock: Schools (3) + Users (8)
- [x] Integración en Container
- [x] Configuración toggle mock/real

### Sprint 2 (Académicos) - ✅ COMPLETADO

- [x] AcademicUnitRepository Mock (13 métodos)
- [x] UnitMembershipRepository Mock (8 métodos)
- [x] SubjectRepository Mock (6 métodos)
- [x] UnitRepository Mock (5 métodos)
- [x] Datos mock: Academic Units (12) + Memberships (5) + Subjects (6) + Units (4)

### Sprint 3 (Completitud) - ✅ COMPLETADO

- [x] MaterialRepository Mock (2 métodos)
- [x] GuardianRepository Mock (13 métodos)
- [x] StatsRepository Mock (1 método con cálculos dinámicos)
- [x] Datos mock: Materials (4) + Guardian Relations (3)
- [x] MockFactory al 100%
- [x] Tests unitarios
- [x] Documentación actualizada

### ✅ IMPLEMENTACIÓN 100% COMPLETA

**9 de 9 repositorios implementados** | **64 métodos totales** | **~2,500 LOC**

## ⚠️ Limitaciones

- **Sin persistencia**: Los datos se pierden al reiniciar la API
- **Sin transacciones**: Cada operación es atómica e independiente
- **Sin locking distribuido**: Solo concurrency control local con RWMutex
- **Capacidad limitada**: Diseñado para desarrollo, no producción

## 🔧 Comandos Útiles

### Desarrollo

```bash
# Activar mocks
export USE_MOCK_REPOSITORIES=true
make run

# Ver logs de startup (confirmar modo)
make run | grep "Usando"
# Output esperado: "✅ Usando MOCK repositories (sin PostgreSQL)"

# Desactivar mocks
export USE_MOCK_REPOSITORIES=false
make run
```

### Verificación

```bash
# Compilar
make build

# Linting
golangci-lint run

# Format
gofmt -w .

# Tests con mocks
USE_MOCK_REPOSITORIES=true go test ./...

# Tests con PostgreSQL
USE_MOCK_REPOSITORIES=false go test ./...
```

## 📚 Referencias

- **Plan de Trabajo**: Ver análisis 360 completo en documentación del proyecto
- **Datos de Seed**: Basados en `edugo-infrastructure/postgres/migrations/testing/`
- **Patrón Factory**: `internal/factory/repository_factory.go`
- **Container DI**: `internal/container/container.go`

## 👥 Contribuir

Para agregar nuevos repositorios mock:

1. Crear datos en `mock/data/entity_name.go`
2. Implementar mock en `mock/repository/entity_repository_mock.go`
3. Actualizar `MockFactory` en `internal/factory/mock_factory.go`
4. Agregar tests de integración

---

**Versión**: Sprints 1, 2 y 3 (COMPLETO - 100%)  
**Última actualización**: 2025-01-29  
**Total Repositorios**: 9/9 (100%)  
**Total Métodos Mock**: 64  
**Total Datos Mock**: 42 registros en 8 entidades
