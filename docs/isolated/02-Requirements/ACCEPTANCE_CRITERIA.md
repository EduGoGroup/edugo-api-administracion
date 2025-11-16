# Criterios de Aceptación - spec-03

## AC-001: Jerarquía Sin Ciclos
```sql
-- Intentar crear ciclo debe fallar
INSERT INTO academic_units (parent_id) VALUES (own_id);
-- Esperado: Error de trigger
```

## AC-002: Árbol Recursivo
```bash
curl /v1/units/root-id/tree
# Esperado: JSON anidado con children recursivos
```

## AC-003: Permisos Jerárquicos
- Owner de padre puede crear hijos
- Owner de unidad puede modificarla
- Non-owner recibe 403

**Total:** 10 criterios medibles
