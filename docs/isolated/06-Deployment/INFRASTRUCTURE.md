# Infrastructure - spec-03
## Docker Compose
```yaml
api-admin:
  image: edugo-api-admin:latest
  ports: ["8081:8081"]
  environment:
    - DB_HOST=postgres
    - PORT=8081
```
## Recursos
- CPU: 1 core
- RAM: 2GB
