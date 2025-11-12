-- =====================================================
-- Script: 01_academic_hierarchy.sql
-- Descripción: Schema de Jerarquía Académica para EduGo
-- Proyecto: edugo-api-administracion
-- Fecha: 12 de Noviembre, 2025
-- =====================================================

-- Habilitar extensión UUID si no está habilitada
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =====================================================
-- TABLA: school
-- Descripción: Escuelas del sistema
-- =====================================================
CREATE TABLE IF NOT EXISTS school (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    address TEXT,
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT school_name_not_empty CHECK (LENGTH(TRIM(name)) > 0),
    CONSTRAINT school_code_not_empty CHECK (LENGTH(TRIM(code)) > 0),
    CONSTRAINT school_email_format CHECK (contact_email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$' OR contact_email IS NULL)
);

CREATE INDEX IF NOT EXISTS idx_school_code ON school(code);
CREATE INDEX IF NOT EXISTS idx_school_created_at ON school(created_at DESC);

CREATE OR REPLACE FUNCTION update_school_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_school_updated_at ON school;
CREATE TRIGGER trigger_school_updated_at
    BEFORE UPDATE ON school
    FOR EACH ROW
    EXECUTE FUNCTION update_school_updated_at();

COMMENT ON TABLE school IS 'Escuelas del sistema EduGo';
COMMENT ON COLUMN school.code IS 'Código único de la escuela (ej: ESC-001)';
COMMENT ON COLUMN school.metadata IS 'Metadata adicional en formato JSON';

-- =====================================================
-- TABLA: academic_unit
-- Descripción: Unidades académicas con jerarquía
-- =====================================================
CREATE TABLE IF NOT EXISTS academic_unit (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_unit_id UUID REFERENCES academic_unit(id) ON DELETE SET NULL,
    school_id UUID NOT NULL REFERENCES school(id) ON DELETE CASCADE,
    unit_type VARCHAR(50) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    code VARCHAR(50),
    description TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    CONSTRAINT academic_unit_type_valid CHECK (
        unit_type IN ('school', 'grade', 'section', 'club', 'department')
    ),
    CONSTRAINT academic_unit_name_not_empty CHECK (LENGTH(TRIM(display_name)) > 0),
    CONSTRAINT academic_unit_no_self_reference CHECK (id != parent_unit_id),
    CONSTRAINT academic_unit_code_unique UNIQUE (school_id, code)
);

CREATE INDEX IF NOT EXISTS idx_academic_unit_school_id ON academic_unit(school_id);
CREATE INDEX IF NOT EXISTS idx_academic_unit_parent_id ON academic_unit(parent_unit_id);
CREATE INDEX IF NOT EXISTS idx_academic_unit_type ON academic_unit(unit_type);
CREATE INDEX IF NOT EXISTS idx_academic_unit_deleted_at ON academic_unit(deleted_at);
CREATE INDEX IF NOT EXISTS idx_academic_unit_school_type ON academic_unit(school_id, unit_type) WHERE deleted_at IS NULL;

CREATE OR REPLACE FUNCTION update_academic_unit_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_academic_unit_updated_at ON academic_unit;
CREATE TRIGGER trigger_academic_unit_updated_at
    BEFORE UPDATE ON academic_unit
    FOR EACH ROW
    EXECUTE FUNCTION update_academic_unit_updated_at();

COMMENT ON TABLE academic_unit IS 'Unidades académicas con estructura jerárquica';
COMMENT ON COLUMN academic_unit.parent_unit_id IS 'Referencia a la unidad padre (para jerarquía)';
COMMENT ON COLUMN academic_unit.unit_type IS 'Tipo: school, grade, section, club, department';
COMMENT ON COLUMN academic_unit.deleted_at IS 'Soft delete: timestamp de eliminación';

-- =====================================================
-- FUNCIÓN: Prevenir ciclos jerárquicos
-- =====================================================
CREATE OR REPLACE FUNCTION prevent_academic_unit_cycles()
RETURNS TRIGGER AS $$
DECLARE
    current_parent_id UUID;
    visited_ids UUID[];
    depth INTEGER := 0;
    max_depth INTEGER := 50;
BEGIN
    IF NEW.parent_unit_id IS NULL THEN
        RETURN NEW;
    END IF;
    
    current_parent_id := NEW.parent_unit_id;
    visited_ids := ARRAY[NEW.id];
    
    WHILE current_parent_id IS NOT NULL AND depth < max_depth LOOP
        IF current_parent_id = ANY(visited_ids) THEN
            RAISE EXCEPTION 'Ciclo detectado en jerarquía: no se puede asignar % como padre de %', 
                NEW.parent_unit_id, NEW.id;
        END IF;
        
        visited_ids := array_append(visited_ids, current_parent_id);
        
        SELECT parent_unit_id INTO current_parent_id
        FROM academic_unit
        WHERE id = current_parent_id;
        
        depth := depth + 1;
    END LOOP;
    
    IF depth >= max_depth THEN
        RAISE EXCEPTION 'Profundidad máxima de jerarquía excedida (máx: %)', max_depth;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_prevent_academic_unit_cycles ON academic_unit;
CREATE TRIGGER trigger_prevent_academic_unit_cycles
    BEFORE INSERT OR UPDATE OF parent_unit_id ON academic_unit
    FOR EACH ROW
    EXECUTE FUNCTION prevent_academic_unit_cycles();

-- =====================================================
-- TABLA: unit_membership
-- Descripción: Relación usuarios-unidades académicas
-- =====================================================
CREATE TABLE IF NOT EXISTS unit_membership (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    unit_id UUID NOT NULL REFERENCES academic_unit(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    role VARCHAR(50) NOT NULL,
    valid_from TIMESTAMP NOT NULL DEFAULT NOW(),
    valid_until TIMESTAMP,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT unit_membership_role_valid CHECK (
        role IN ('student', 'teacher', 'coordinator', 'admin', 'assistant')
    ),
    CONSTRAINT unit_membership_dates_valid CHECK (
        valid_until IS NULL OR valid_until > valid_from
    ),
    CONSTRAINT unit_membership_unique UNIQUE (unit_id, user_id, valid_from)
);

CREATE INDEX IF NOT EXISTS idx_unit_membership_unit_id ON unit_membership(unit_id);
CREATE INDEX IF NOT EXISTS idx_unit_membership_user_id ON unit_membership(user_id);
CREATE INDEX IF NOT EXISTS idx_unit_membership_role ON unit_membership(role);
CREATE INDEX IF NOT EXISTS idx_unit_membership_valid_dates ON unit_membership(valid_from, valid_until);

CREATE OR REPLACE FUNCTION update_unit_membership_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_unit_membership_updated_at ON unit_membership;
CREATE TRIGGER trigger_unit_membership_updated_at
    BEFORE UPDATE ON unit_membership
    FOR EACH ROW
    EXECUTE FUNCTION update_unit_membership_updated_at();

COMMENT ON TABLE unit_membership IS 'Membresías de usuarios en unidades académicas';
COMMENT ON COLUMN unit_membership.role IS 'Rol: student, teacher, coordinator, admin, assistant';
COMMENT ON COLUMN unit_membership.valid_from IS 'Fecha de inicio de la membresía';
COMMENT ON COLUMN unit_membership.valid_until IS 'Fecha de fin (NULL = indefinido)';

-- =====================================================
-- VISTA: v_unit_tree
-- Descripción: Vista con árbol jerárquico (CTE recursivo)
-- =====================================================
CREATE OR REPLACE VIEW v_unit_tree AS
WITH RECURSIVE unit_hierarchy AS (
    SELECT 
        id,
        parent_unit_id,
        school_id,
        unit_type,
        display_name,
        code,
        description,
        1 AS depth,
        ARRAY[id] AS path,
        display_name::TEXT AS full_path
    FROM academic_unit
    WHERE parent_unit_id IS NULL
      AND deleted_at IS NULL
    
    UNION ALL
    
    SELECT 
        au.id,
        au.parent_unit_id,
        au.school_id,
        au.unit_type,
        au.display_name,
        au.code,
        au.description,
        uh.depth + 1,
        uh.path || au.id,
        (uh.full_path || ' > ' || au.display_name)::TEXT
    FROM academic_unit au
    INNER JOIN unit_hierarchy uh ON au.parent_unit_id = uh.id
    WHERE au.deleted_at IS NULL
)
SELECT 
    uh.*,
    s.name AS school_name,
    s.code AS school_code
FROM unit_hierarchy uh
LEFT JOIN school s ON uh.school_id = s.id
ORDER BY uh.school_id, uh.path;

COMMENT ON VIEW v_unit_tree IS 'Vista con árbol jerárquico completo de unidades académicas';

-- =====================================================
-- VISTA: v_active_memberships
-- Descripción: Vista de membresías activas
-- =====================================================
CREATE OR REPLACE VIEW v_active_memberships AS
SELECT 
    um.id,
    um.unit_id,
    um.user_id,
    um.role,
    um.valid_from,
    um.valid_until,
    um.metadata,
    au.display_name AS unit_name,
    au.unit_type,
    au.school_id,
    s.name AS school_name
FROM unit_membership um
INNER JOIN academic_unit au ON um.unit_id = au.id
INNER JOIN school s ON au.school_id = s.id
WHERE (um.valid_until IS NULL OR um.valid_until > NOW())
  AND au.deleted_at IS NULL
ORDER BY s.name, au.display_name, um.role, um.valid_from DESC;

COMMENT ON VIEW v_active_memberships IS 'Vista de membresías activas con información de unidades';
