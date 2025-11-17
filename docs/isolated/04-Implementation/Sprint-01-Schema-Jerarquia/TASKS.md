# Tareas Sprint 01

## TASK-01-001: Crear Migración SQL
**Estimación:** 4h

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion/scripts/postgresql/01_academic_hierarchy.sql`

```sql
BEGIN;

CREATE TABLE school (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE academic_unit (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    school_id UUID NOT NULL REFERENCES school(id) ON DELETE CASCADE,
    parent_unit_id UUID REFERENCES academic_unit(id) ON DELETE CASCADE,
    unit_type VARCHAR(20) NOT NULL CHECK (unit_type IN ('grade','section','club','department')),
    display_name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    depth INTEGER DEFAULT 0,
    path TEXT, -- Materialized path: /school/grade/section
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(school_id, code),
    CHECK (id != parent_unit_id)
);

CREATE INDEX idx_unit_parent ON academic_unit(parent_unit_id);
CREATE INDEX idx_unit_school ON academic_unit(school_id);
CREATE INDEX idx_unit_path ON academic_unit(path);

CREATE TABLE unit_membership (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    unit_id UUID NOT NULL REFERENCES academic_unit(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL CHECK (role IN ('student','teacher','owner','admin')),
    joined_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(unit_id, user_id)
);

CREATE INDEX idx_membership_unit ON unit_membership(unit_id);
CREATE INDEX idx_membership_user ON unit_membership(user_id);

-- Trigger para actualizar depth automáticamente
CREATE OR REPLACE FUNCTION update_unit_depth() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.parent_unit_id IS NULL THEN
        NEW.depth := 0;
    ELSE
        SELECT depth + 1 INTO NEW.depth
        FROM academic_unit WHERE id = NEW.parent_unit_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_depth
    BEFORE INSERT OR UPDATE ON academic_unit
    FOR EACH ROW EXECUTE FUNCTION update_unit_depth();

COMMIT;
```

#### Criterios
- [ ] 3 tablas creadas
- [ ] Índices optimizados
- [ ] Trigger actualiza depth
- [ ] Foreign keys con CASCADE

**Tiempo:** 4h
