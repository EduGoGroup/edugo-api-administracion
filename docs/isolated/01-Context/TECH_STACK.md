# TECH STACK - API Admin

## Resumen Ejecutivo

| Layer | Tecnología | Versión | Propósito |
|-------|-----------|---------|----------|
| **Language** | Go | 1.21+ | Backend compilado |
| **Framework** | Gin | v1.9+ | Web framework HTTP |
| **ORM** | GORM | v1.25+ | Abstracción PostgreSQL |
| **Primary DB** | PostgreSQL | 15+ | Datos relacionales + CTEs recursivas |
| **Config Management** | Viper | Latest | Multi-environment |
| **Authentication** | JWT (shared) | Custom | Token-based auth |
| **Containerization** | Docker | 20.10+ | Container runtime |

---

## Stack Detallado

### Capa 1: Aplicación (Go)

Similar a api-mobile, pero con énfasis en queries complejas y CTEs.

```go
package main

func main() {
    // Inicializar
    logger.Init(config.Logger)
    db := database.InitPostgres(config.Database)
    
    // Crear router
    router := gin.New()
    router.Use(auth.ValidateToken())
    
    // Registrar handlers
    schoolHandler := handlers.NewSchoolHandler(db)
    hierarchyHandler := handlers.NewHierarchyHandler(db)
    
    v1 := router.Group("/api/v1")
    {
        v1.POST("/schools", schoolHandler.Create)
        v1.GET("/schools", schoolHandler.List)
        v1.GET("/schools/:id/hierarchy", hierarchyHandler.GetTree)
    }
    
    router.Run(":8081")
}
```

---

### Capa 2: Web Framework (Gin)

Mismo que api-mobile.

---

### Capa 3: ORM (GORM)

```go
type School struct {
    ID                int64              `gorm:"primaryKey"`
    Name              string
    AcademicUnits     []AcademicUnit     `gorm:"foreignKey:SchoolID"`
    Teachers          []Teacher          `gorm:"foreignKey:SchoolID"`
    Students          []Student          `gorm:"foreignKey:SchoolID"`
}

type AcademicUnit struct {
    ID        int64
    SchoolID  int64
    ParentID  *int64
    Type      string
    Name      string
    Children  []AcademicUnit `gorm:"foreignKey:ParentID"`
    Parent    *AcademicUnit  `gorm:"foreignKey:ParentID"`
}

// Cargar con relaciones
var school School
db.Preload("AcademicUnits").First(&school, id)
```

---

### Capa 4: Base de Datos (PostgreSQL + CTEs)

**Característica especial:** Queries recursivas para manejar árbol académico.

```sql
-- Obtener árbol completo
WITH RECURSIVE hierarchy AS (
    SELECT id, parent_id, name, 0 as level
    FROM academic_units
    WHERE school_id = $1 AND parent_id IS NULL
    
    UNION ALL
    
    SELECT au.id, au.parent_id, au.name, h.level + 1
    FROM academic_units au
    INNER JOIN hierarchy h ON au.parent_id = h.id
)
SELECT * FROM hierarchy ORDER BY level, name;
```

**Índices críticos:**
```sql
CREATE INDEX idx_au_school_parent ON academic_units(school_id, parent_id);
```

---

### Capa 5: Configuración (Viper)

```go
viper.SetDefault("database.host", "localhost")
viper.SetDefault("database.port", 5432)
viper.SetDefault("api.port", 8081)
```

---

## Características Especiales

### 1. Queries Recursivas (CTEs)

API Admin utiliza ampliamente Common Table Expressions recursivas:

```go
// En repository
func (r *AcademicUnitRepository) GetHierarchy(ctx context.Context, unitID int64) ([]HierarchyNode, error) {
    const query = `
    WITH RECURSIVE tree AS (
        SELECT id, parent_id, name, 0 as depth
        FROM academic_units
        WHERE id = ?
        
        UNION ALL
        
        SELECT au.id, au.parent_id, au.name, t.depth + 1
        FROM academic_units au
        INNER JOIN tree t ON au.parent_id = t.id
    )
    SELECT * FROM tree ORDER BY depth, name
    `
    
    var nodes []HierarchyNode
    result := r.db.WithContext(ctx).Raw(query, unitID).Scan(&nodes)
    return nodes, result.Error
}
```

### 2. Multi-tenancy (por school_id)

Todas las queries incluyen filtro de school_id:

```go
// Scope automático
db.Where("school_id = ?", schoolID).Find(&units)
```

### 3. Validaciones de Árbol

```go
// No permitir ciclos
func (r *AcademicUnitRepository) ValidateParent(ctx context.Context, unitID, parentID int64) error {
    // Verificar que parentID no es descendiente de unitID
    const query = `
    WITH RECURSIVE descendants AS (...)
    SELECT COUNT(*) FROM descendants WHERE id = ?
    `
    
    var count int64
    r.db.Raw(query, parentID, unitID).Scan(&count)
    if count > 0 {
        return errors.New("would create cycle")
    }
    return nil
}
```

---

## Comparativa con api-mobile

| Aspecto | API Mobile | API Admin |
|--------|-----------|----------|
| Queries | Simples CRUD | Complejas (CTEs) |
| Datos | Evaluaciones | Jerarquía |
| Relacionales | Pocas relaciones | Árbol complejo |
| Performance crítica | Medium | HIGH (índices vitales) |
| Concurrencia | Alta (muchos users) | Media (administradores) |
| Caché | Útil | Muy útil |

---

## Performance Considerations

### Índices Críticos
```sql
-- OBLIGATORIOS para performance
CREATE INDEX idx_au_school_parent ON academic_units(school_id, parent_id);
CREATE INDEX idx_teachers_school ON teachers(school_id);
CREATE INDEX idx_students_school ON students(school_id);
CREATE INDEX idx_memberships_unit ON memberships(academic_unit_id);
```

### Límites de Profundidad
```go
const MAX_HIERARCHY_DEPTH = 10

func (s *HierarchyService) ValidateDepth(ctx context.Context, parentID int64) error {
    // Contar profundidad actual
    depth := 0
    current := parentID
    
    for current != 0 && depth < MAX_HIERARCHY_DEPTH {
        parent := s.repo.GetParent(current)
        current = parent
        depth++
    }
    
    if depth >= MAX_HIERARCHY_DEPTH {
        return errors.New("hierarchy too deep")
    }
    
    return nil
}
```

---

## Benchmark Esperado

**Operación típica:** Obtener árbol de 100 unidades académicas
- Query time: ~50-100ms
- Memory: ~5MB
- Con índices: ~10-20ms

---

## Versionamiento

Mismo que api-mobile (basado en SHARED)
