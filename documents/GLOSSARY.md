# 📖 Glosario de Términos

> Definiciones de conceptos y términos utilizados en EduGo API Administración

## 🏛️ Términos del Dominio Educativo

### School (Escuela)
Institución educativa principal. Es la entidad raíz de toda la jerarquía académica.
- **Ejemplo**: "Colegio San Martín", "Instituto Técnico Nacional"
- **Identificador**: UUID único + código alfanumérico único
- **Relaciones**: Contiene múltiples Academic Units

### Academic Unit (Unidad Académica)
Cualquier nivel organizacional dentro de una escuela. Forma una estructura jerárquica (árbol).

| Tipo | Descripción | Ejemplo | Padre Típico |
|------|-------------|---------|--------------|
| `school` | Nivel raíz (la escuela misma) | - | null |
| `grade` | Grado o año escolar | "1° Primaria", "5° Secundaria" | school |
| `section` | División de un grado | "Sección A", "Turno Mañana" | grade |
| `department` | Área académica | "Departamento de Matemáticas" | school |
| `club` | Actividad extracurricular | "Club de Ajedrez" | school/department |

### Membership (Membresía)
Relación entre un usuario y una unidad académica con un rol específico.
- **Características**: Tiene fechas de validez, puede ser activa o expirada
- **Un usuario puede tener múltiples membresías** en diferentes unidades

### Membership Roles (Roles de Membresía)

| Rol | Descripción | Permisos Típicos |
|-----|-------------|------------------|
| `director` | Director de la institución | Acceso total a la escuela |
| `coordinator` | Coordinador académico | Gestión de unidades asignadas |
| `teacher` | Profesor | Acceso a sus unidades/materias |
| `assistant` | Asistente | Soporte en unidades específicas |
| `student` | Estudiante | Acceso a contenido de sus unidades |
| `observer` | Observador | Solo lectura |

---

## 👥 Términos de Usuarios

### User (Usuario)
Persona registrada en el sistema. Puede tener un rol del sistema y múltiples membresías.

### System Roles (Roles del Sistema)

| Rol | Descripción | Alcance |
|-----|-------------|---------|
| `super_admin` | Administrador global | Todo el sistema |
| `school_admin` | Admin de escuela | Una escuela específica |
| `teacher` | Profesor | Sus unidades asignadas |
| `student` | Estudiante | Sus unidades matriculadas |
| `guardian` | Padre/Tutor | Información de sus hijos |

### Guardian (Tutor)
Usuario con rol de padre/tutor que tiene relación con uno o más estudiantes.
- **Relaciones**: `father`, `mother`, `guardian`, `other`
- **Primary Guardian**: Tutor principal para notificaciones

---

## 🔐 Términos de Autenticación

### JWT (JSON Web Token)
Token de autenticación firmado digitalmente que contiene información del usuario.

```
Header.Payload.Signature
```

### Access Token
Token de corta duración (15 minutos) para autenticar requests.
- **Uso**: Header `Authorization: Bearer {token}`
- **Renovación**: Usar refresh token cuando expira

### Refresh Token
Token de larga duración (7 días) para obtener nuevos access tokens.
- **Uso**: Solo en endpoint `/v1/auth/refresh`
- **Seguridad**: Debe almacenarse de forma segura

### Claims
Datos contenidos en el JWT:
- `sub`: User ID
- `email`: Email del usuario
- `role`: Rol del sistema
- `iss`: Issuer (edugo-central)
- `exp`: Timestamp de expiración

### Token Blacklist
Lista de tokens revocados (por logout) que ya no son válidos aunque no hayan expirado.

### Issuer
Identificador del servicio que emitió el token. En EduGo siempre es `edugo-central`.

---

## 🏗️ Términos de Arquitectura

### Clean Architecture
Patrón arquitectónico que separa el código en capas concéntricas:
- **Domain**: Entidades y reglas de negocio (centro)
- **Application**: Casos de uso y orquestación
- **Infrastructure**: Frameworks, DB, HTTP (exterior)

### Repository Pattern
Abstracción que encapsula el acceso a datos. Define una interfaz (contrato) que puede tener múltiples implementaciones (PostgreSQL, Mock, etc.).

```go
type SchoolRepository interface {
    Create(ctx context.Context, school *entities.School) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error)
    // ...
}
```

### Service Layer
Capa que contiene la lógica de negocio. Orquesta repositorios y aplica reglas.

### Handler
Componente que maneja requests HTTP. Parsea input, llama servicios, formatea output.

### DTO (Data Transfer Object)
Objeto para transferir datos entre capas. Diferente de las entidades de dominio.
- **Request DTO**: Datos de entrada (CreateSchoolRequest)
- **Response DTO**: Datos de salida (SchoolResponse)

### Dependency Injection (DI)
Patrón donde las dependencias se pasan (inyectan) a los componentes en lugar de crearlas internamente.

```go
// El handler recibe el service, no lo crea
func NewSchoolHandler(service SchoolService) *SchoolHandler
```

### Container
Componente central que crea e inyecta todas las dependencias de la aplicación.

---

## 🗄️ Términos de Base de Datos

### Soft Delete
Técnica donde los registros no se eliminan físicamente, solo se marca un timestamp `deleted_at`.
- **Ventaja**: Permite recuperación y auditoría
- **Implementación**: GORM filtra automáticamente registros con `deleted_at != null`

### UUID
Identificador único universal (Universal Unique Identifier). Formato: `550e8400-e29b-41d4-a716-446655440000`

### JSONB
Tipo de datos de PostgreSQL para almacenar JSON de forma binaria indexable.
- **Uso en EduGo**: Campo `metadata` en varias entidades

### CTE (Common Table Expression)
Expresión de tabla común en SQL. Usado para queries recursivas como jerarquías.

```sql
WITH RECURSIVE hierarchy AS (
    SELECT * FROM academic_unit WHERE id = :id
    UNION ALL
    SELECT au.* FROM academic_unit au
    JOIN hierarchy h ON au.id = h.parent_unit_id
)
SELECT * FROM hierarchy;
```

---

## 📡 Términos de API

### REST
Representational State Transfer. Estilo arquitectónico para APIs web.

### Endpoint
URL específica que maneja un tipo de operación.
- `GET /v1/schools` - Listar escuelas
- `POST /v1/schools` - Crear escuela

### HTTP Status Codes

| Código | Significado | Uso |
|--------|-------------|-----|
| 200 | OK | Operación exitosa |
| 201 | Created | Recurso creado |
| 204 | No Content | Operación exitosa sin respuesta |
| 400 | Bad Request | Error de validación |
| 401 | Unauthorized | No autenticado |
| 403 | Forbidden | Sin permisos |
| 404 | Not Found | Recurso no existe |
| 409 | Conflict | Conflicto (ej: duplicado) |
| 429 | Too Many Requests | Rate limit excedido |
| 500 | Internal Server Error | Error del servidor |

### Rate Limiting
Limitación del número de requests por unidad de tiempo para prevenir abuso.

### Swagger/OpenAPI
Especificación para documentar APIs REST. Genera documentación interactiva.

---

## 🧪 Términos de Testing

### Unit Test
Test que prueba una unidad de código en aislamiento (función, método).

### Integration Test
Test que prueba la interacción entre componentes (API + DB real).

### Testcontainers
Librería que levanta contenedores Docker para tests de integración.

### Mock
Implementación falsa de una interfaz para testing.

```go
type MockSchoolRepository struct {
    CreateFunc func(ctx context.Context, school *entities.School) error
}
```

### Test Coverage
Porcentaje de código cubierto por tests.

---

## 🔄 Términos de Operaciones

### Health Check
Endpoint que verifica si el servicio está funcionando.
```
GET /health → {"status": "healthy"}
```

### Graceful Shutdown
Proceso de apagado que espera a que las operaciones en curso terminen antes de cerrar.

### Environment Variables
Variables de configuración cargadas desde el sistema operativo o archivo `.env`.

### Middleware
Función que se ejecuta antes/después de cada request. Ejemplos: autenticación, logging.

---

## 📦 Términos de Ecosistema EduGo

### edugo-shared
Paquetes compartidos entre servicios EduGo:
- `auth`: JWT manager
- `bootstrap`: Inicialización
- `common`: Tipos comunes
- `logger`: Logging estructurado

### edugo-infrastructure
Paquetes de infraestructura:
- `postgres`: Entidades y conexión DB

### API Mobile (edugo-api-mobile)
API para aplicaciones móviles. Consume tokens del servicio de auth centralizado.

### Workers
Servicios de procesamiento en background. También usan auth centralizado.
