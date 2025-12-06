# 🔧 Mejoras y Deuda Técnica

> Código deprecado, refactorizaciones pendientes y mejoras identificadas

## 📋 Índice

| Documento | Contenido |
|-----------|-----------|
| [DEPRECATED_CODE.md](./DEPRECATED_CODE.md) | Código marcado como deprecado a eliminar |
| [REFACTORING.md](./REFACTORING.md) | Refactorizaciones pendientes |
| [CODE_SMELLS.md](./CODE_SMELLS.md) | Malas prácticas identificadas |
| [TODO_LIST.md](./TODO_LIST.md) | Lista de TODOs encontrados en el código |

---

## 🎯 Resumen Ejecutivo

### Prioridad Alta 🔴

| Item | Archivo | Descripción | Impacto |
|------|---------|-------------|---------|
| Legacy Handlers | `cmd/legacy_handlers.go` | Handlers deprecados sin uso | Código muerto |
| Valores hardcodeados | `school_service.go` | Country, tier, limits hardcodeados | Configurabilidad |
| Función incompleta | `unit_membership_service.go:174` | `ListMembershipsByRole` no filtra por rol | Bug funcional |

### Prioridad Media 🟡

| Item | Archivo | Descripción | Impacto |
|------|---------|-------------|---------|
| Código repetitivo | Handlers | Error handling duplicado | Mantenibilidad |
| Validación de roles | Services | Roles hardcodeados en arrays | Extensibilidad |
| Parámetro sin usar | Services | `activeOnly` no se usa en queries | Funcionalidad incompleta |

### Prioridad Baja 🟢

| Item | Archivo | Descripción | Impacto |
|------|---------|-------------|---------|
| TODOs pendientes | Varios | Comentarios TODO sin resolver | Documentación |
| Tests faltantes | Varios | Cobertura incompleta | Calidad |

---

## 📊 Métricas de Deuda Técnica

```
Total de items identificados:    23
├── Código deprecado:            6
├── Refactorizaciones:           8
├── Code smells:                 5
└── TODOs pendientes:            4

Estimación de esfuerzo total:    ~20-30 horas de desarrollo
```

---

## ✅ Acciones Recomendadas

### Sprint Actual

1. **Eliminar `legacy_handlers.go`** - Ya marcado para v0.6.0
2. **Implementar `ListMembershipsByRole` correctamente** - Es un bug
3. **Extraer roles válidos a constantes/config**

### Próximo Sprint

1. **Crear middleware genérico de error handling**
2. **Mover valores hardcodeados a configuración**
3. **Usar parámetro `activeOnly` en queries**

### Backlog

1. **Aumentar cobertura de tests a 80%+**
2. **Implementar cache con Redis**
3. **Agregar métricas Prometheus**

---

## 🔄 Proceso de Mejora

```
1. Identificar issue          → Agregar a este documento
2. Crear ticket en backlog    → Referencia este doc
3. Implementar fix            → PR con tests
4. Actualizar documentación   → Marcar como resuelto
5. Eliminar de este documento → En siguiente release
```

---

## 📝 Historial de Mejoras Completadas

| Fecha | Item | Descripción | PR |
|-------|------|-------------|-----|
| 2025-11-20 | Auth centralizado | Migración a auth unificado | #45 |
| 2025-11-15 | Clean Architecture | Refactor a capas limpias | #42 |
| 2025-11-01 | Bootstrap | Migración a shared/bootstrap | #38 |
