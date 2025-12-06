# Validación Sprint 01
```bash
# Ejecutar migración
psql -f scripts/postgresql/01_academic_hierarchy.sql

# Verificar tablas
psql -c "\d school"
psql -c "\d academic_unit"
psql -c "\d unit_membership"

# Test de trigger depth
psql -c "INSERT INTO school VALUES (gen_uuid_v7(), 'Test', 'TEST');"
psql -c "INSERT INTO academic_unit (school_id, display_name, code, unit_type) VALUES ((SELECT id FROM school WHERE code='TEST'), 'Unit', 'U1', 'grade');"
psql -c "SELECT depth FROM academic_unit WHERE code='U1';"
# Esperado: 0
```
## Criterios
- [ ] Tablas creadas
- [ ] Trigger funciona
- [ ] Índices creados
