# 🚀 START HERE - API Admin (Jerarquía Académica)

## ⭐ PROYECTO COMPLETADO ✅

**Estado:** ✅ COMPLETADO (v0.2.0)  
**Fecha finalización:** 12 de Noviembre, 2025

**Bienvenido a la documentación COMPLETA y AUTÓNOMA de edugo-api-administracion.**

Esta carpeta contiene la documentación del proyecto YA IMPLEMENTADO de jerarquía académica.

**📍 Documentación oficial:** `/Users/jhoanmedina/source/EduGo/Analisys/docs/specs/api-admin-jerarquia/`

---

## 📍 ¿Qué es edugo-api-administracion?

**API REST** para gestión administrativa de EduGo (instituciones, estructuras académicas, usuarios).

**Puerto:** 8081  
**Tecnología:** Go + Gin + GORM + PostgreSQL  
**Arquitectura:** Clean Architecture (Hexagonal)

### Funcionalidades Principales
- ✅ **Gestión de Escuelas** (CRUD de instituciones)
- ✅ **Estructura Académica Jerárquica** (árbol de unidades académicas)
- ✅ **Memberships** (asignación de usuarios a roles y unidades)
- ✅ **Consultas Recursivas** (búsqueda en árboles jerárquicos)
- ✅ **Autenticación Administrativa** (JWT con roles específicos)

---

## 🎯 ¿Qué Se Implementó? (COMPLETADO)

**Sistema Completo de Jerarquía Académica:**

1. **Escuelas (Schools)**
   - CRUD básico de instituciones
   - Metadatos (nombre, código, ubicación)
   - Logo y configuración institucional

2. **Unidades Académicas (Academic Units)**
   - Estructura árbol jerárquico (parent-child)
   - Tipos: Facultad → Departamento → Carrera → Programa
   - Consultas recursivas (ancestor, descendant)
   - Búsqueda rápida por código

3. **Memberships (Asignaciones)**
   - Asignar usuario a unidad académica
   - Roles: DIRECTOR, DOCENTE, COORDINADOR, ADMIN
   - Permisos basados en rol y unidad
   - Historial de asignaciones

4. **Reportes Administrativos**
   - Estructura académica completa
   - Usuarios por unidad
   - Estadísticas de memberships

---

## 📂 Estructura de Esta Carpeta

```
api-admin/
│
├── START_HERE.md                ⭐ Este archivo - LEER PRIMERO
├── EXECUTION_PLAN.md            Plan paso a paso de ejecución
│
├── 01-Context/                  Contexto del proyecto
│   ├── PROJECT_OVERVIEW.md      Overview detallado
│   ├── ECOSYSTEM_CONTEXT.md     Cómo encaja en el ecosistema
│   ├── DEPENDENCIES.md          Qué necesita de otros proyectos
│   └── TECH_STACK.md            Stack tecnológico
│
├── 02-Requirements/             Requisitos funcionales y técnicos
│   ├── PRD.md                   Product Requirements Document
│   ├── FUNCTIONAL_SPECS.md      Especificaciones funcionales
│   ├── TECHNICAL_SPECS.md       Especificaciones técnicas
│   └── ACCEPTANCE_CRITERIA.md   Criterios de aceptación
│
├── 03-Design/                   Diseño arquitectónico
│   ├── ARCHITECTURE.md          Arquitectura Clean detallada
│   ├── DATA_MODEL.md            Modelo de datos con jerarquía
│   ├── RECURSIVE_QUERIES.md     Estrategia de consultas recursivas
│   ├── API_CONTRACTS.md         Contratos de API (OpenAPI)
│   └── SECURITY_DESIGN.md       Diseño de autenticación y autorización
│
├── 04-Implementation/           Implementación (6 sprints)
│   ├── Sprint-01-Schema-BD/     Schema PostgreSQL y jerarquía
│   ├── Sprint-02-Dominio/       Entities, Value Objects, Interfaces
│   ├── Sprint-03-Repositorios/  Implementación de repositorios
│   ├── Sprint-04-Services-API/  Services y endpoints REST
│   ├── Sprint-05-Testing/       Tests unitarios e integración
│   └── Sprint-06-CI-CD/         CI/CD y deployment
│
├── 05-Testing/                  Estrategia de testing
│   ├── TEST_STRATEGY.md
│   ├── TEST_CASES.md
│   └── COVERAGE_REPORT.md
│
├── 06-Deployment/               Deployment y monitoreo
│   ├── DEPLOYMENT_GUIDE.md
│   ├── INFRASTRUCTURE.md
│   └── MONITORING.md
│
└── PROGRESS.json                Tracking de progreso (JSON)
```

---

## 🚦 Flujo de Inicio Rápido

### Paso 1: Leer Contexto (15 min)
```bash
# Entender qué es este proyecto y cómo encaja
cat 01-Context/PROJECT_OVERVIEW.md
cat 01-Context/ECOSYSTEM_CONTEXT.md
cat 01-Context/DEPENDENCIES.md
```

### Paso 2: Revisar Requisitos (30 min)
```bash
# Entender QUÉ vamos a construir
cat 02-Requirements/PRD.md
cat 02-Requirements/FUNCTIONAL_SPECS.md
cat 02-Requirements/ACCEPTANCE_CRITERIA.md
```

### Paso 3: Estudiar Arquitectura (45 min)
```bash
# Entender CÓMO lo vamos a construir
cat 03-Design/ARCHITECTURE.md
cat 03-Design/DATA_MODEL.md
cat 03-Design/RECURSIVE_QUERIES.md
```

### Paso 4: Ejecutar Plan (Ver EXECUTION_PLAN.md)
```bash
# Plan detallado de implementación
cat EXECUTION_PLAN.md
```

### Paso 5: Implementar Sprint por Sprint (18 días estimados)
```bash
cd 04-Implementation/Sprint-01-Schema-BD/
cat README.md
cat TASKS.md
# ... ejecutar tareas ...
# Repetir para cada sprint
```

---

## 🔗 Dependencias Externas

Este proyecto **NECESITA** de otros componentes del ecosistema:

### 1. edugo-infrastructure v0.1.1
**Versión usada:** v0.1.1  
**Qué se usó:**
- `database/migrations/001_create_users.up.sql`
- `database/migrations/002_create_schools.up.sql`
- `database/migrations/005_create_academic_hierarchy.up.sql`
- `database/TABLE_OWNERSHIP.md` - Documenta ownership de tablas

**Estado:** ✅ Implementado y funcionando

### 2. edugo-shared v0.7.0
**Versión usada:** v0.7.0 (FROZEN)  
**Módulos usados:**
- `config` - Configuración multi-ambiente
- `database/postgres` - Conexiones PostgreSQL
- `auth` - JWT y autenticación
- `logger` - Logging estructurado
- `bootstrap` - Dependency injection (creado en FASE 0.1 de este proyecto)

**Estado:** ✅ Funcionando perfectamente

### 3. PostgreSQL 15+
**Uso:** Base de datos principal (jerarquía académica)  
**Tablas implementadas:** ✅
- `schools` (escuelas/instituciones)
- `academic_units` (estructura jerárquica)
- `unit_memberships` (asignaciones usuario-unidad-rol)

**Características especiales:**
- Soporte de CTEs (Common Table Expressions) para recursión
- Índices BTREE en foreign keys
- Índices HASH para búsquedas por código

### 4. Base de Datos Existente
**Tablas previas requeridas:**
- `users` (usuarios sistema)
- `roles` (roles globales)

**Cambios:** Agregar columnas `created_at`, `updated_at`, `deleted_at`

---

## ⚙️ Configuración Requerida

### Variables de Entorno
```bash
# PostgreSQL
DATABASE_URL=postgres://user:pass@localhost:5432/edugo_dev?sslmode=disable

# Auth
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# Server
PORT=8081
ENVIRONMENT=local  # local, dev, qa, prod
LOG_LEVEL=debug

# Características (feature flags)
ENABLE_HIERARCHY_CACHE=true
HIERARCHY_CACHE_TTL=3600  # 1 hora
```

### Prerequisitos de Sistema
```bash
# Go 1.21+
go version

# PostgreSQL 15+ con soporte de CTEs
psql --version

# Docker (para desarrollo)
docker --version
```

---

## ✅ Fases Completadas (TODAS)

**FASE 0.1-0.3:** Bootstrap Compartido ✅  
- shared/bootstrap creado
- api-mobile y worker migrados
- 2,667 LOC en shared, -937 LOC en mobile

**FASE 1:** Modernización ✅  
- Clean Architecture implementada
- PRs #12, #13 merged

**FASE 2:** Schema BD ✅  
- 3 tablas + constraints + seeds
- PR #15 merged

**FASE 3:** Dominio ✅  
- 3 entities, 8 value objects
- PR #16 merged

**FASE 4:** Services ✅  
- Services + DTOs + Repositories
- PR #17 merged

**FASE 5:** API REST ✅  
- 15+ endpoints funcionales
- PR #18 merged

**FASE 6:** Testing ✅  
- Suite completa >80% coverage
- PR #19 merged

**FASE 7:** CI/CD ✅  
- GitHub Actions workflows
- PR #20 merged

**Release:** v0.2.0 publicado

---

## ✅ Checklist Pre-Implementación

Antes de comenzar Sprint 01, verifica:

### Ambiente de Desarrollo
- [ ] Go 1.21+ instalado
- [ ] PostgreSQL 15+ corriendo (con soporte CTE)
- [ ] Repositorio edugo-api-administracion clonado
- [ ] Rama feature creada: `git checkout -b feature/academic-hierarchy`

### Dependencias
- [ ] edugo-shared v1.3.0 publicado en GitHub
- [ ] Tabla `users` existe en PostgreSQL
- [ ] PostgreSQL versión >= 15 (para CTEs)

### Configuración
- [ ] Archivo `.env.local` creado con variables necesarias
- [ ] Conexión a PostgreSQL verificada: `psql -U user -d edugo_dev`
- [ ] Soporte de CTEs verificado: `psql -U user -d edugo_dev -c "WITH RECURSIVE..."`

### Opcional
- [ ] Al menos 1 usuario de prueba en tabla `users`
- [ ] Roles básicos definidos en tabla `roles`

---

## 🎯 Resultado Esperado

Al completar los 6 sprints, tendrás:

### Funcionalidades
- ✅ API REST completa de jerarquía académica
- ✅ 12+ endpoints funcionando
- ✅ Consultas recursivas eficientes
- ✅ CRUD completo de escuelas y unidades
- ✅ Sistema de memberships y roles
- ✅ Reportes de estructura académica

### Calidad
- ✅ Cobertura de tests >85%
- ✅ Tests de integración con Testcontainers
- ✅ Documentación Swagger actualizada
- ✅ CI/CD funcionando (GitHub Actions)

### Arquitectura
- ✅ Clean Architecture implementada
- ✅ Domain Layer independiente
- ✅ Repositorios con interfaces
- ✅ Código mantenible y testeable
- ✅ Consultas optimizadas con caché

---

## 📞 Soporte y Recursos

### Dentro de Esta Carpeta
- **Dudas de arquitectura:** `03-Design/ARCHITECTURE.md`
- **Dudas de jerarquía:** `03-Design/DATA_MODEL.md`
- **Dudas de consultas recursivas:** `03-Design/RECURSIVE_QUERIES.md`
- **Dudas de requisitos:** `02-Requirements/`
- **Dudas de implementación:** `04-Implementation/Sprint-XX/TASKS.md`
- **Dudas de testing:** `05-Testing/TEST_STRATEGY.md`

### Contexto del Ecosistema
- **Cómo encaja este proyecto:** `01-Context/ECOSYSTEM_CONTEXT.md`
- **Qué depende de qué:** `01-Context/DEPENDENCIES.md`
- **Stack tecnológico:** `01-Context/TECH_STACK.md`

---

## 🚀 Comenzar AHORA

```bash
# 1. Lee el overview del proyecto
cat 01-Context/PROJECT_OVERVIEW.md

# 2. Lee el plan de ejecución
cat EXECUTION_PLAN.md

# 3. Inicia Sprint 01
cd 04-Implementation/Sprint-01-Schema-BD/
cat README.md
cat TASKS.md

# 4. Ejecuta las tareas paso a paso
# ... sigue las instrucciones de TASKS.md
```

---

**Última actualización:** 15 de Noviembre, 2025  
**Generado con:** Claude Code  
**Proyecto:** edugo-api-administracion - Jerarquía Académica  
**Tipo de documentación:** Aislada y autónoma

---

## 🎓 Filosofía de Esta Documentación

> **"Todo lo que necesitas está aquí. No necesitas buscar en archivos externos. Esta carpeta es autónoma."**

**Si encuentras que falta algo, es un bug en la documentación. Repórtalo.**

---

¡Éxito en tu implementación! 🚀
