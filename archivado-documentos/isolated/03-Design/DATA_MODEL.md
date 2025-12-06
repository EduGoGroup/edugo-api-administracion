# Modelo de Datos - spec-03

## Tablas PostgreSQL

### school
```sql
CREATE TABLE school (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### academic_unit
```sql
CREATE TABLE academic_unit (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    school_id UUID NOT NULL REFERENCES school(id),
    parent_unit_id UUID REFERENCES academic_unit(id),
    unit_type VARCHAR(20) CHECK (unit_type IN ('grade','section','club')),
    display_name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    depth INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(school_id, code),
    
    -- Trigger para prevenir ciclos
    CONSTRAINT no_self_reference CHECK (id != parent_unit_id)
);

CREATE INDEX idx_unit_parent ON academic_unit(parent_unit_id);
CREATE INDEX idx_unit_school ON academic_unit(school_id);
```

### unit_membership
```sql
CREATE TABLE unit_membership (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    unit_id UUID NOT NULL REFERENCES academic_unit(id),
    user_id UUID NOT NULL REFERENCES users(id),
    role VARCHAR(20) CHECK (role IN ('student','teacher','owner')),
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(unit_id, user_id)
);
```

## Query Recursiva para Árbol
```sql
WITH RECURSIVE unit_tree AS (
    SELECT id, parent_unit_id, display_name, depth
    FROM academic_unit
    WHERE id = $1
    
    UNION ALL
    
    SELECT u.id, u.parent_unit_id, u.display_name, u.depth
    FROM academic_unit u
    INNER JOIN unit_tree t ON u.parent_unit_id = t.id
)
SELECT * FROM unit_tree;
```
