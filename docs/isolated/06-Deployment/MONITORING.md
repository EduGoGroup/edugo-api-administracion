# Monitoring - spec-03
## Métricas
- GET /units/:id/tree latency: <500ms
- CRUD operations: <200ms
- Error rate: <1%
## Prometheus
```go
UnitsCreated = prometheus.NewCounter(...)
TreeQueryDuration = prometheus.NewHistogram(...)
```
