# Changelog

Todos los cambios notables en edugo-api-administracion serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.7.0] - 2025-11-25

### Tipo de Release: patch

- release: merge dev to main - Auth centralizada Sprint 1 (#51)
- chore: actualizar infraestructura a v0.10.1

---

## [0.6.3] - 2025-11-22

### Tipo de Release: patch

- fix: agregar BindEnv para AUTH_JWT_SECRET en configuración
- style: aplicar gofmt a archivos de migración
- fix(main): eliminar rutas deprecated undefined

---

## [0.6.1] - 2025-11-22

### Tipo de Release: patch

- hotfix: eliminar rutas deprecated undefined en main.go

---

## [0.6.0] - 2025-11-22

### Tipo de Release: patch

- fix(test): eliminar import no usado en academic_unit_service_test
- feat: Migrar de DDD a Infrastructure Entities
- test: agregar tests unitarios para services (cobertura +2.8%)
- fix(cicd): agregar permisos para acceso a repos privados
- fix(cicd): agregar GOPRIVATE en step de descarga de dependencias
- fix(lint): corregir errores de errcheck en repositorios
- fix(cicd): remover replace local para permitir descarga desde GitHub
- feat(migration): eliminar DDD completamente y cleanup final (FASE 4)
- feat(migration): migrar services restantes y DTOs a infrastructure (FASE 3)
- feat(migration): migrar 6 entidades restantes a infrastructure (FASE 2)
- feat(migration): migrar User de DDD a infrastructure entities (FASE 1)
- feat(cicd): migrar workflows a reusables centralizados
- docs(sprint-4): agregar sección PR a main (FASE 3 extendida)
- docs(sprint-4): agregar sección FASE 3 a lecciones aprendidas
- docs(sprint-4): agregar lecciones aprendidas de api-mobile
- docs(sprint-2): completar SPRINT-2 - documentación final FASE 3
- fix(lint): corregir todos los defer rows.Close() restantes
- fix(lint): corregir 4 errcheck detectados por golangci-lint v2
- fix(ci): actualizar golangci-lint-action v6 -> v7
- fix(ci): usar golangci-lint v2.6.2 para soportar Go 1.25


---

## [0.5.1] - 2025-11-19

### Tipo de Release: patch

- chore: release v0.5.0

---

## [0.5.0] - 2025-11-19

### Tipo de Release: minor



---

## [0.4.4] - 2025-11-18

### Tipo de Release: patch

- fix: copiar archivos de configuración en Dockerfile (#37)

---

## [0.4.3] - 2025-11-18

### Tipo de Release: patch

- release: Sprint-04 - All tests passing (17/17) (#36)

---

## [0.4.2] - 2025-11-18

### Tipo de Release: patch

- release: Sprint-04 - HTTP REST API with ltree support (Fixed) (#35)

---

## [0.4.0] - 2025-11-18

### Tipo de Release: patch

- release: Sprint-04 - HTTP REST API with ltree support (#34)
- docs: clarify Sprint-04 must branch from dev (not main)
- docs: add Sprint-04 workflow prompts (Phase 1 Web + Phase 2 Local)
- chore: Sync dev to main - Sprint-03 Repositorios ltree (#32)
- docs: add phase 1 (web) and phase 2 (local) workflow prompts
- feat(domain): Sprint-02 - Tree traversal, domain services, and test optimization (#29)
- Revert "feat(domain): Sprint-02 - Tree traversal, domain services, and test optimization (#28)"
- feat(domain): Sprint-02 - Tree traversal, domain services, and test optimization (#28)
- feat: Migración a Infrastructure Centralizado v0.7.1 (#27)

---

## [0.3.1] - 2025-11-17

### Tipo de Release: patch



---

## [0.2.0] - 2025-11-13

### Tipo de Release: minor

- release: v0.2.0 - Sistema de Jerarquía Académica Completo (#21)

---

## [0.1.2] - 2025-11-12

### Tipo de Release: patch

- chore: actualizar shared a v0.4.0 + modernización completada (#14)
- docs: documentar GitHub App y actualizar a v2.1.4
- feat: implementar GitHub App Token para sincronización automática

---

## [0.1.1] - 2025-11-03

### Tipo de Release: patch

- fix: corregir workflow sync-main-to-dev (#4)

---

## [0.1.0] - 2025-11-01

### Tipo de Release: minor

- feat: CI/CD optimizado y Copilot instructions (#3)

---

## [Unreleased]

## [0.1.0] - 2025-11-01

### Added
- Sistema GitFlow profesional implementado
- Workflows de CI/CD automatizados:
  - CI Pipeline con tests y validaciones
  - Tests con cobertura y servicios de infraestructura
  - Manual Release workflow (TODO-EN-UNO) para control total de releases
  - Docker only workflow para builds manuales
  - Release automático con versionado semántico
  - Sincronización automática main ↔ dev
- GitHub Copilot custom instructions en español
- Migración a edugo-shared con arquitectura modular
- Submódulos: common, logger
- .gitignore completo para Go
- Documentación completa de workflows

### Changed
- Actualizado a Go 1.25.3
- Versionado corregido a v0.x.x (proyecto en desarrollo)
- Eliminado auto-version.yml (reemplazado por manual-release.yml)

### Fixed
- Corrección de errores de linter (errcheck)
- Permisos de GitHub Container Registry configurados

[Unreleased]: https://github.com/EduGoGroup/edugo-api-administracion/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/EduGoGroup/edugo-api-administracion/releases/tag/v0.1.0
