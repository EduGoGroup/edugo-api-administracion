# Tareas Sprint 06
## TASK-06-001: GitHub Actions
```yaml
# .github/workflows/ci.yml
jobs:
  test:
    services:
      postgres: ...
    steps:
      - run: go test ./... -v
```
**Tiempo:** 3h
