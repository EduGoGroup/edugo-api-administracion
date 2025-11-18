# Changelog

Todos los cambios notables en edugo-api-administracion serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
