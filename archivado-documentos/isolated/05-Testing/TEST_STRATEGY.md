# Estrategia Testing - spec-03
## Pirámide: 70% unit, 20% integration, 10% E2E
## Coverage: >80% global
## Tests Críticos:
- Query recursiva árbol
- Prevención ciclos
- Permisos jerárquicos
```bash
go test ./... -cover
```
