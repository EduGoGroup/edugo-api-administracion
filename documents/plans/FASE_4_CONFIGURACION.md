# 🟡 FASE 4: Mejoras de Configuración

**Prioridad**: Media  
**Estimación**: 4 horas  
**Rama**: `feat/fase-4-configuracion`

---

## Preparación Git

```bash
git checkout dev
git pull origin dev
git checkout -b feat/fase-4-configuracion
```

---

## 4.1 Agregar campos City y Country al DTO de School

### Problema Actual
Country hardcodeado como "CO" y City siempre nil.

### Ubicación del Problema
```
internal/application/service/school_service.go:68-85
```

### Código Actual (Problemático)
```go
school := &entities.School{
    // ...
    City:    nil,   // TODO: agregar cuando se agregue al DTO
    Country: "CO",  // Hardcoded
    // ...
}
```

### Tareas
1. Agregar `City` y `Country` a `CreateSchoolRequest`
2. Agregar `City` y `Country` a `UpdateSchoolRequest`
3. Actualizar mapeo en el servicio
4. Actualizar documentación Swagger

### Modificar DTO
```
internal/application/dto/school_dto.go
```

### Código
```go
type CreateSchoolRequest struct {
    Name         string                 `json:"name" validate:"required,min=3"`
    Code         string                 `json:"code" validate:"required,min=3"`
    Address      string                 `json:"address"`
    City         string                 `json:"city"`
    Country      string                 `json:"country"`
    ContactEmail string                 `json:"contact_email" validate:"omitempty,email"`
    ContactPhone string                 `json:"contact_phone"`
    Metadata     map[string]interface{} `json:"metadata"`
}

type UpdateSchoolRequest struct {
    Name         *string                `json:"name" validate:"omitempty,min=3"`
    Address      *string                `json:"address"`
    City         *string                `json:"city"`
    Country      *string                `json:"country"`
    ContactEmail *string                `json:"contact_email" validate:"omitempty,email"`
    ContactPhone *string                `json:"contact_phone"`
    Metadata     map[string]interface{} `json:"metadata"`
}
```

### Esfuerzo
1 hora

---

## 4.2 Mover valores de suscripción a configuración

### Problema Actual
`SubscriptionTier`, `MaxTeachers`, `MaxStudents` hardcodeados en el servicio.

### Código Actual (Problemático)
```go
school := &entities.School{
    // ...
    SubscriptionTier: "free",  // Hardcoded
    MaxTeachers:      50,      // Hardcoded
    MaxStudents:      500,     // Hardcoded
    // ...
}
```

### Tareas
1. Agregar sección `defaults.school` en configuración
2. Crear struct `SchoolDefaults`
3. Inyectar configuración en `SchoolService`
4. Usar defaults o valores del request

### Modificar Configuración
```
config/config.yaml
```

### Agregar Sección
```yaml
defaults:
  school:
    country: "CO"
    subscription_tier: "free"
    max_teachers: 50
    max_students: 500
```

### Crear Struct de Configuración
```
internal/config/school_defaults.go
```

### Código
```go
package config

// SchoolDefaults contiene los valores por defecto para escuelas
type SchoolDefaults struct {
    Country          string `yaml:"country" env:"DEFAULT_SCHOOL_COUNTRY" env-default:"CO"`
    SubscriptionTier string `yaml:"subscription_tier" env:"DEFAULT_SCHOOL_TIER" env-default:"free"`
    MaxTeachers      int    `yaml:"max_teachers" env:"DEFAULT_MAX_TEACHERS" env-default:"50"`
    MaxStudents      int    `yaml:"max_students" env:"DEFAULT_MAX_STUDENTS" env-default:"500"`
}

// Defaults contiene todas las configuraciones de valores por defecto
type Defaults struct {
    School SchoolDefaults `yaml:"school"`
}
```

### Modificar Service
```go
type schoolService struct {
    schoolRepo repository.SchoolRepository
    logger     logger.Logger
    defaults   config.SchoolDefaults
}

func NewSchoolService(
    repo repository.SchoolRepository, 
    logger logger.Logger,
    defaults config.SchoolDefaults,
) SchoolService {
    return &schoolService{
        schoolRepo: repo,
        logger:     logger,
        defaults:   defaults,
    }
}

func (s *schoolService) CreateSchool(ctx context.Context, req dto.CreateSchoolRequest) (*dto.SchoolResponse, error) {
    // Usar defaults si no se proporciona valor
    country := req.Country
    if country == "" {
        country = s.defaults.Country
    }

    tier := req.SubscriptionTier
    if tier == "" {
        tier = s.defaults.SubscriptionTier
    }

    school := &entities.School{
        // ...
        Country:          country,
        SubscriptionTier: tier,
        MaxTeachers:      s.defaults.MaxTeachers,
        MaxStudents:      s.defaults.MaxStudents,
        // ...
    }
}
```

### Actualizar Container
Modificar la inyección de dependencias para pasar los defaults al service.

### Esfuerzo
3 horas

---

## Documentación a Actualizar

Al completar esta fase, actualizar:

- `documents/improvements/CODE_SMELLS.md` - Eliminar sección 1 (valores hardcodeados)
- `documents/improvements/TODO_LIST.md` - Eliminar TODOs de City/Country/Subscription
- `documents/API.md` - Actualizar documentación de endpoint de School con nuevos campos
- `documents/SETUP.md` - Documentar nuevas opciones de configuración en la sección de variables de entorno

---

## Finalización

```bash
git add .
git commit -m "feat: agregar campos City/Country y mover defaults a configuración"
git push origin feat/fase-4-configuracion
```

### Crear PR a dev con:
- Título: `feat: agregar campos City/Country y mover defaults a configuración`
- Descripción: Fase 4 del plan de mejoras - Configuración

---

## Checklist

- [ ] `City` y `Country` agregados a DTOs
- [ ] Sección `defaults.school` agregada a config.yaml
- [ ] Struct `SchoolDefaults` creado
- [ ] `SchoolService` modificado para usar defaults
- [ ] Container actualizado para inyectar defaults
- [ ] Swagger actualizado
- [ ] Tests actualizados
- [ ] Documentación actualizada
- [ ] PR creado a dev
