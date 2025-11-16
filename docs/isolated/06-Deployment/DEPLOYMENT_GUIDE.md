# Deployment - spec-03
## Pasos
1. Ejecutar migración 01_academic_hierarchy.sql
2. Build: `go build -o bin/api-admin ./cmd/api`
3. Deploy con Docker o systemd
4. Puerto 8081
5. Health check: `curl http://localhost:8081/health`
```bash
docker run -p 8081:8081 edugo-api-admin
```
