# CI/CD Workflows

**Versión**: 2.0 (Simplificado)  
**Última actualización**: 9 de noviembre de 2025

---

## 🚀 Quick Start

### Para Desarrolladores

```bash
# 1. Crear feature branch
git checkout -b feature/nueva-funcionalidad

# 2. Desarrollar y commit
git add .
git commit -m "feat: nueva funcionalidad"

# 3. Push y crear PR a dev
git push origin feature/nueva-funcionalidad
# → pr-to-dev.yml se ejecuta automáticamente (~2-3 min)

# 4. Después de merge, crear PR de dev a main
# → pr-to-main.yml se ejecuta automáticamente (~3-4 min)
```

### Para Releases

```
1. Ir a: Actions → Manual Release → Run workflow
2. Ingresar versión (ej: 0.1.0)
3. Seleccionar tipo: patch/minor/major
4. Click "Run workflow"
→ Release completo en ~10-15 minutos
```

---

## 📋 Workflows Activos (5)

| Workflow | Trigger | Propósito | Tiempo |
|----------|---------|-----------|--------|
| **pr-to-dev.yml** | Auto (PR a dev) | Tests rápidos | 2-3 min |
| **pr-to-main.yml** | Auto (PR a main) | Tests completos | 3-4 min |
| **manual-release.yml** | Manual | Release completo | 10-15 min |
| **test.yml** | Manual | Tests on-demand | Variable |
| **sync-main-to-dev.yml** | Auto (push a main) | Sync branches | 30s |

---

## 📚 Documentación

Toda la documentación está en [`docs/`](docs/):

- **[WORKFLOWS_INDEX.md](docs/WORKFLOWS_INDEX.md)** - 📋 Índice completo de workflows
- **[CI_CD_STRATEGY.md](docs/CI_CD_STRATEGY.md)** - 🎯 Estrategia de CI/CD
- **[WORKFLOW_DIAGRAM.md](docs/WORKFLOW_DIAGRAM.md)** - 📊 Diagramas visuales
- **[SIMPLIFICATION_PLAN.md](docs/SIMPLIFICATION_PLAN.md)** - 🔧 Plan de simplificación
- **[TESTING_STRATEGY.md](docs/TESTING_STRATEGY.md)** - 🧪 Estrategia de testing
- **[TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)** - 🐛 Resolución de problemas

---

## 🎯 Flujo Simplificado

```
feature → dev → main → release
   ↓       ↓      ↓       ↓
  PR     PR    Merge  Manual
   ↓       ↓      ↓       ↓
 2-3min 3-4min  30s   10-15min
```

---

## ✅ Mejoras Implementadas

- ✅ **54% menos workflows** (11 → 5)
- ✅ **Sin duplicación** de código
- ✅ **81.5% más rápido** tests de integración
- ✅ **Documentación organizada** en `docs/`
- ✅ **Todo on-demand** excepto sync

---

## 🔗 Links Rápidos

- [Ver Workflows en GitHub](../../actions)
- [Crear Release Manual](../../actions/workflows/manual-release.yml)
- [Ejecutar Tests Manual](../../actions/workflows/test.yml)
- [Ver Documentación Completa](docs/WORKFLOWS_INDEX.md)

---

**¿Problemas?** Ver [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)

