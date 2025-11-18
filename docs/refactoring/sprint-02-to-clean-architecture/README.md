# Refactor: De DDD Rico a Clean Architecture Estricta

**Proyecto:** edugo-api-administracion  
**Fecha Inicio:** 2025-11-17  
**Estimación:** 25-30 horas  
**Prioridad:** Media  
**Estado:** 🟡 Planificación

---

## 📋 Índice de Documentos

1. **[README.md](README.md)** - Este archivo (overview del refactor)
2. **[IMPACT_ANALYSIS.md](IMPACT_ANALYSIS.md)** - Análisis de impacto y riesgos
3. **[WORK_PLAN.md](WORK_PLAN.md)** - Plan de trabajo detallado con fases
4. **[TARGET_ARCHITECTURE.md](TARGET_ARCHITECTURE.md)** - Arquitectura objetivo
5. **[VALIDATION_CHECKLIST.md](VALIDATION_CHECKLIST.md)** - Checklist de validación
6. **[PROGRESS_TRACKING.md](PROGRESS_TRACKING.md)** - Seguimiento de progreso

---

## 🎯 Objetivo

Migrar de un enfoque **DDD Rico** (lógica en entities) a **Clean Architecture Estricta** (lógica en domain services), manteniendo la funcionalidad actual mientras mejoramos la separación de responsabilidades.

### Estado Actual (DDD Rico)

```
┌─────────────────────────────────────────┐
│         Entity (AcademicUnit)           │
│  ┌───────────────────────────────────┐  │
│  │  - id, name, type, children       │  │
│  │  - SetParent()                    │  │
│  │  - AddChild()                     │  │
│  │  - GetAllDescendants()            │  │
│  │  - GetDepth()                     │  │
│  │  - Todas las validaciones         │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

### Estado Objetivo (Clean Architecture)

```
┌─────────────────────────────────────────┐
│    Entity (AcademicUnit) - Anemic      │
│  ┌───────────────────────────────────┐  │
│  │  - id, name, type, children       │  │
│  │  - Getters/Setters simples        │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
                    ▲
                    │ usa
                    │
┌─────────────────────────────────────────┐
│   Domain Service (AcademicUnitService)  │
│  ┌───────────────────────────────────┐  │
│  │  - SetParent(unit, parent)        │  │
│  │  - AddChild(parent, child)        │  │
│  │  - GetAllDescendants(unit)        │  │
│  │  - GetDepth(unit)                 │  │
│  │  - Todas las validaciones         │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

---

## 🎪 Motivación

### Problemas Actuales
1. **Entities con demasiada responsabilidad**: Mezclan estado y comportamiento
2. **Difícil de testear**: Mock de entities complejas
3. **No es arquitectura limpia pura**: Va contra principios SOLID estrictos
4. **Confusión en `.coverignore`**: Se asumió que entities eran "solo structs"

### Beneficios Esperados
1. ✅ **Separación clara**: Entities = datos, Services = lógica
2. ✅ **Más testeable**: Services independientes, entities simples
3. ✅ **Arquitectura limpia by-the-book**: Cumple con Uncle Bob
4. ✅ **Escalabilidad**: Fácil agregar nuevos services

---

## 📊 Métricas de Éxito

| Métrica | Antes | Objetivo |
|---------|-------|----------|
| Líneas de código en Entity | ~400 | ~150 |
| Líneas de código en Service | 0 | ~300 |
| Tests de Entity | 656 líneas | ~200 líneas |
| Tests de Service | 0 | ~500 líneas |
| Cobertura total | 13.2% | >35% |
| Archivos modificados | - | ~35 |

---

## ⚠️ Riesgos Principales

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Romper invariantes | Alta | Alto | Tests exhaustivos antes/después |
| Tiempo > estimación | Media | Medio | Trabajo por fases incrementales |
| Introducir bugs | Media | Alto | Code review estricto + tests |
| Conflictos con otros PRs | Baja | Medio | Comunicar con equipo |

---

## 📅 Timeline

```
Semana 1: Planificación y setup
  ├─ Día 1-2: Documentación y análisis
  └─ Día 3-5: Fase 1 (Domain Services base)

Semana 2: Migración core
  ├─ Día 1-3: Fase 2 (Migrar entities)
  └─ Día 4-5: Fase 3 (Tests)

Semana 3: Integración y validación
  ├─ Día 1-2: Fase 4 (Repositorios y app layer)
  ├─ Día 3-4: Fase 5 (Validación completa)
  └─ Día 5: PR y review
```

---

## 👥 Stakeholders

- **Ejecutor**: Claude Code
- **Revisor**: @medinatello
- **Aprobador**: Tech Lead del proyecto
- **Impactados**: Desarrolladores que usen estas entities

---

## 📚 Referencias

1. [Clean Architecture - Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
2. [DDD vs Anemic Domain Model](https://martinfowler.com/bliki/AnemicDomainModel.html)
3. [Go Clean Architecture Examples](https://github.com/bxcodec/go-clean-arch)
4. [Effective Go - Best Practices](https://go.dev/doc/effective_go)

---

## 🚦 Estado de Documentos

- [x] README.md (este archivo)
- [ ] IMPACT_ANALYSIS.md
- [ ] WORK_PLAN.md
- [ ] TARGET_ARCHITECTURE.md
- [ ] VALIDATION_CHECKLIST.md
- [ ] PROGRESS_TRACKING.md

---

**Última actualización:** 2025-11-17  
**Próximo paso:** Completar documentación de análisis de impacto
