# 🔄 Workflow Templates - Ejecución en 2 Fases

**Fecha:** 16 de Noviembre, 2025  
**Versión:** 1.0.0  
**Propósito:** Templates genéricos y reutilizables para workflow de 2 fases (Web + Local)

---

## 🎯 ¿Qué es esto?

Esta carpeta contiene **templates genéricos** para implementar workflows de 2 fases en cualquier proyecto:

- **Fase 1 (Claude Code Web):** Implementación con stubs/mocks para recursos externos
- **Fase 2 (Claude Code Local):** Implementación real con Docker/DB, CI/CD, merge

---

## 📦 Templates Incluidos

### 1. WORKFLOW_ORCHESTRATION.md
Sistema completo de orquestación de 2 fases con:
- Flujo detallado Fase 1 (Web)
- Flujo detallado Fase 2 (Local)
- Reglas de trabajo desatendido
- Manejo de errores y timeouts
- Monitoreo de CI/CD

### 2. TRACKING_SYSTEM.md
Sistema de tracking con PROGRESS.json:
- Estados de sprint
- Timestamps de fases
- Métricas de progreso

### 3. PHASE2_BRIDGE_TEMPLATE.md
Template para documento puente entre fases:
- Lista de stubs creados
- Código stub vs. código real requerido
- Validaciones pendientes
- Checklist para Fase 2

### 4. PROGRESS_TEMPLATE.json
Template de archivo de tracking JSON:
- Estructura de sprints
- Estados y timestamps
- Métricas de proyecto

### 5. scripts/
Scripts de automatización:
- `update-progress.sh` - Actualizar PROGRESS.json
- `recover.sh` - Recuperación ante interrupciones
- `daily-report.sh` - Reporte diario de progreso

---

## 🚀 Cómo Usar Estos Templates

### Paso 1: Copiar Templates a Tu Proyecto

```bash
# Ir a tu proyecto en 00-Projects-Isolated
cd /path/to/your/project/docs/isolated/

# Copiar archivos base
cp /path/to/workflow-templates/WORKFLOW_ORCHESTRATION.md ./
cp /path/to/workflow-templates/TRACKING_SYSTEM.md ./
cp /path/to/workflow-templates/PROGRESS_TEMPLATE.json ./PROGRESS.json

# Copiar scripts
mkdir -p scripts
cp -r /path/to/workflow-templates/scripts/* ./scripts/
```

### Paso 2: Adaptar PROGRESS.json a Tu Proyecto

Editar `PROGRESS.json` con los sprints específicos de tu proyecto:

```json
{
  "project": "tu-proyecto-nombre",
  "version": "1.0.0",
  "current_sprint": "Sprint-01-Nombre",
  "sprints": {
    "Sprint-01-Nombre": {
      "name": "Descripción del Sprint",
      "status": "pending",
      "estimated_hours": 8,
      "tasks": {
        "TASK-001": {
          "name": "Nombre de la tarea",
          "status": "pending"
        }
      }
    }
  }
}
```

### Paso 3: Crear PHASE2_BRIDGE.md por Sprint

Para cada sprint, crear:
```
04-Implementation/Sprint-XX-Nombre/PHASE2_BRIDGE.md
```

Usar template: `PHASE2_BRIDGE_TEMPLATE.md`

---

## 📋 Estructura Recomendada en Cada Proyecto

```
tu-proyecto/
├── docs/
│   └── isolated/
│       ├── WORKFLOW_ORCHESTRATION.md     ← Copiado de template
│       ├── TRACKING_SYSTEM.md             ← Copiado de template  
│       ├── PROGRESS.json                  ← Adaptado del template
│       │
│       ├── 04-Implementation/
│       │   ├── Sprint-01-.../
│       │   │   ├── README.md
│       │   │   ├── TASKS.md
│       │   │   ├── DEPENDENCIES.md
│       │   │   ├── VALIDATION.md
│       │   │   ├── PHASE2_BRIDGE.md       ← Generado en Fase 1
│       │   │   └── EXECUTION_REPORT.md    ← Generado en Fase 2
│       │   │
│       │   └── Sprint-02-.../
│       │       └── [misma estructura]
│       │
│       └── scripts/
│           ├── update-progress.sh         ← Copiado de template
│           ├── recover.sh                 ← Copiado de template
│           └── daily-report.sh            ← Copiado de template
```

---

## ✅ Beneficios del Workflow

1. **Ejecución desatendida** en Claude Code Web (sin Docker/DB)
2. **Continuación local** con recursos reales y CI/CD
3. **Recuperación automática** ante interrupciones
4. **Tracking detallado** de progreso por sprint
5. **CI/CD validado** antes de merge a dev
6. **Code review** de Copilot atendido automáticamente

---

## 🎯 Proyectos que Usan Estos Templates

- `edugo-api-mobile` - Sistema de Evaluaciones
- `edugo-worker` - Procesamiento de Materiales
- `edugo-admin-api` - API de Administración
- *(Agregar otros proyectos aquí)*

---

## 📚 Documentación Adicional

- **Flujo completo Fase 1 y 2:** Ver `WORKFLOW_ORCHESTRATION.md`
- **Sistema de tracking:** Ver `TRACKING_SYSTEM.md`
- **Template de bridge:** Ver `PHASE2_BRIDGE_TEMPLATE.md`

---

## 🔄 Versionado de Templates

**Versión actual:** 1.0.0

### Changelog
- **1.0.0** (2025-11-16) - Templates iniciales extraídos de edugo-api-mobile

---

## 🤝 Contribuir

Si mejoras estos templates, por favor:
1. Actualizar versión en este README
2. Documentar cambios en Changelog
3. Notificar a proyectos que ya usan los templates

---

**Última actualización:** 16 de Noviembre, 2025  
**Mantenedor:** EduGo Team  
**Licencia:** Interno EduGo

---

## 🎓 Filosofía

> **"Reutiliza el proceso, no el código. Estos templates son el proceso."**

Los templates permiten:
- ✅ Consistencia entre proyectos
- ✅ Onboarding rápido de nuevos proyectos
- ✅ Mejores prácticas documentadas
- ✅ Workflow probado y refinado

---

¡Éxito en tu proyecto! 🚀
