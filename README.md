# EduGo API Administración

API REST para operaciones administrativas y CRUD en la plataforma EduGo.

## Descripción

Esta API maneja:
- Gestión de usuarios (crear, editar, eliminar)
- Gestión de jerarquía académica (escuelas, unidades)
- Gestión de materias
- Moderación de contenidos
- Estadísticas globales del sistema

## Tecnología

- Go 1.21+ + Gin + Swagger
- Puerto: `8081`

## Instalación

```bash
go mod download
swag init -g cmd/main.go -o docs
go run cmd/main.go
```

## Endpoints

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/v1/users` | Crear usuario |
| PATCH | `/v1/users/:id` | Actualizar usuario |
| DELETE | `/v1/users/:id` | Eliminar usuario |
| POST | `/v1/schools` | Crear escuela |
| POST | `/v1/units` | Crear unidad académica |
| PATCH | `/v1/units/:id` | Actualizar unidad |
| POST | `/v1/units/:id/members` | Asignar membresía |
| POST | `/v1/subjects` | Crear materia |
| DELETE | `/v1/materials/:id` | Eliminar material |
| GET | `/v1/stats/global` | Estadísticas globales |

## Swagger

`http://localhost:8081/swagger/index.html`

## Estado: Código base con datos MOCK

Implementar para producción:
- Validación de rol admin en middleware
- Lógica CRUD real con PostgreSQL
- Auditoría en `audit_log`
- Validaciones de negocio (ej: prevenir jerarquía circular)
