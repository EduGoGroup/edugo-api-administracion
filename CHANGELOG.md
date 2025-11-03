# Changelog

Todos los cambios notables en edugo-api-administracion serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
