-- Migration: 012_extend_for_admin_api
-- Description: Extender schema para soportar api-admin (jerarquía académica)
-- Dependencies: 004_create_memberships.up.sql
-- Date: 2025-11-17
-- Owner: infrastructure (compartido)

BEGIN;

-- =====================================================
-- 1. Extender academic_units para soportar jerarquía
-- =====================================================

-- 1.1 Agregar parent_unit_id para estructura jerárquica
ALTER TABLE academic_units
    ADD COLUMN IF NOT EXISTS parent_unit_id UUID REFERENCES academic_units(id) ON DELETE SET NULL;

-- 1.2 Agregar índice para parent_unit_id
CREATE INDEX IF NOT EXISTS idx_academic_units_parent ON academic_units(parent_unit_id);

-- 1.3 Agregar constraint para prevenir auto-referencia
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'academic_units_no_self_reference'
          AND conrelid = 'academic_units'::regclass
    ) THEN
        ALTER TABLE academic_units
            ADD CONSTRAINT academic_units_no_self_reference CHECK (id != parent_unit_id);
    END IF;
END $$;

-- 1.4 Extender tipos permitidos (agregar school, club, department)
ALTER TABLE academic_units DROP CONSTRAINT IF EXISTS academic_units_type_check;
ALTER TABLE academic_units ADD CONSTRAINT academic_units_type_check
    CHECK (type IN ('school', 'grade', 'class', 'section', 'club', 'department'));

-- 1.5 Hacer academic_year nullable (no todos los proyectos lo usan)
ALTER TABLE academic_units ALTER COLUMN academic_year DROP NOT NULL;
ALTER TABLE academic_units ALTER COLUMN academic_year SET DEFAULT 0;

-- 1.6 Agregar metadata y description
ALTER TABLE academic_units ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'::jsonb;
ALTER TABLE academic_units ADD COLUMN IF NOT EXISTS description TEXT;

COMMENT ON COLUMN academic_units.parent_unit_id IS 'Unidad padre (para jerarquía: Facultad → Departamento → Carrera)';
COMMENT ON COLUMN academic_units.metadata IS 'Metadata adicional en formato JSON (extensible)';
COMMENT ON COLUMN academic_units.description IS 'Descripción de la unidad académica';
COMMENT ON COLUMN academic_units.academic_year IS 'Año académico (0 = sin año específico)';

-- =====================================================
-- 2. Extender schools con metadata
-- =====================================================

ALTER TABLE schools ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'::jsonb;

COMMENT ON COLUMN schools.metadata IS 'Metadata adicional en formato JSON (logo, config, etc.)';

-- =====================================================
-- 3. Extender memberships con roles administrativos
-- =====================================================

-- 3.1 Extender roles permitidos (agregar coordinator, admin, assistant)
ALTER TABLE memberships DROP CONSTRAINT IF EXISTS memberships_role_check;
ALTER TABLE memberships ADD CONSTRAINT memberships_role_check
    CHECK (role IN ('teacher', 'student', 'guardian', 'coordinator', 'admin', 'assistant'));

-- 3.2 Agregar metadata
ALTER TABLE memberships ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}'::jsonb;

COMMENT ON COLUMN memberships.metadata IS 'Metadata adicional (permisos específicos, historial, etc.)';

-- =====================================================
-- 4. Crear función para prevenir ciclos en jerarquía
-- =====================================================

CREATE OR REPLACE FUNCTION prevent_academic_unit_cycles()
RETURNS TRIGGER AS $$
DECLARE
    current_parent_id UUID;
    visited_ids UUID[];
    depth INTEGER := 0;
    max_depth INTEGER := 50;
BEGIN
    -- Si no hay parent, no hay problema
    IF NEW.parent_unit_id IS NULL THEN
        RETURN NEW;
    END IF;

    current_parent_id := NEW.parent_unit_id;
    visited_ids := ARRAY[]::UUID[];

    -- Agregar el ID actual si no es NULL (en INSERT ya existe)
    IF NEW.id IS NOT NULL THEN
        visited_ids := array_append(visited_ids, NEW.id);
    END IF;

    -- Recorrer hacia arriba en la jerarquía
    WHILE current_parent_id IS NOT NULL AND depth < max_depth LOOP
        -- Detectar ciclo
        IF current_parent_id = ANY(visited_ids) THEN
            RAISE EXCEPTION 'Ciclo detectado en jerarquía: no se puede asignar % como padre de %',
                NEW.parent_unit_id, NEW.id;
        END IF;

        visited_ids := array_append(visited_ids, current_parent_id);

        -- Obtener el siguiente padre
        SELECT parent_unit_id INTO current_parent_id
        FROM academic_units
        WHERE id = current_parent_id;

        depth := depth + 1;
    END LOOP;

    -- Validar profundidad máxima
    IF depth >= max_depth THEN
        RAISE EXCEPTION 'Profundidad máxima de jerarquía excedida (máx: %)', max_depth;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 4.1 Crear trigger para prevenir ciclos
DROP TRIGGER IF EXISTS trigger_prevent_academic_unit_cycles ON academic_units;
CREATE TRIGGER trigger_prevent_academic_unit_cycles
    BEFORE INSERT OR UPDATE OF parent_unit_id ON academic_units
    FOR EACH ROW
    EXECUTE FUNCTION prevent_academic_unit_cycles();

COMMENT ON FUNCTION prevent_academic_unit_cycles() IS 'Previene ciclos en jerarquía de academic_units';

-- =====================================================
-- 5. Crear vista para árbol jerárquico (CTE recursivo)
-- =====================================================

CREATE OR REPLACE VIEW v_academic_unit_tree AS
WITH RECURSIVE unit_hierarchy AS (
    -- Caso base: unidades raíz (sin padre)
    SELECT
        id,
        parent_unit_id,
        school_id,
        name,
        code,
        type,
        level,
        academic_year,
        1 AS depth,
        ARRAY[id] AS path,
        name::TEXT AS full_path
    FROM academic_units
    WHERE parent_unit_id IS NULL
      AND deleted_at IS NULL

    UNION ALL

    -- Caso recursivo: hijos de cada unidad
    SELECT
        au.id,
        au.parent_unit_id,
        au.school_id,
        au.name,
        au.code,
        au.type,
        au.level,
        au.academic_year,
        uh.depth + 1,
        uh.path || au.id,
        (uh.full_path || ' > ' || au.name)::TEXT
    FROM academic_units au
    INNER JOIN unit_hierarchy uh ON au.parent_unit_id = uh.id
    WHERE au.deleted_at IS NULL
)
SELECT
    uh.*,
    s.name AS school_name,
    s.code AS school_code
FROM unit_hierarchy uh
LEFT JOIN schools s ON uh.school_id = s.id
ORDER BY uh.school_id, uh.path;

COMMENT ON VIEW v_academic_unit_tree IS 'Vista con árbol jerárquico completo de unidades académicas';

-- =====================================================
-- 6. Crear vista de memberships activas (mejorada)
-- =====================================================

DROP VIEW IF EXISTS v_active_memberships;
CREATE OR REPLACE VIEW v_active_memberships AS
SELECT
    m.id,
    m.user_id,
    m.school_id,
    m.academic_unit_id,
    m.role,
    m.is_active,
    m.enrolled_at,
    m.withdrawn_at,
    m.metadata,
    au.name AS unit_name,
    au.type AS unit_type,
    au.academic_year,
    s.name AS school_name,
    s.code AS school_code
FROM memberships m
INNER JOIN schools s ON m.school_id = s.id
LEFT JOIN academic_units au ON m.academic_unit_id = au.id
WHERE m.is_active = true
  AND (au.deleted_at IS NULL OR au.id IS NULL)
ORDER BY s.name, au.name, m.role, m.enrolled_at DESC;

COMMENT ON VIEW v_active_memberships IS 'Vista de memberships activas con información completa';

COMMIT;
