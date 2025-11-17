-- Migration: 012_extend_for_admin_api (ROLLBACK)
-- Description: Revertir extensiones para api-admin
-- Date: 2025-11-17

BEGIN;

-- =====================================================
-- 1. Eliminar vistas
-- =====================================================

DROP VIEW IF EXISTS v_active_memberships;
DROP VIEW IF EXISTS v_academic_unit_tree;

-- =====================================================
-- 2. Eliminar función y trigger de ciclos
-- =====================================================

DROP TRIGGER IF EXISTS trigger_prevent_academic_unit_cycles ON academic_units;
DROP FUNCTION IF EXISTS prevent_academic_unit_cycles();

-- =====================================================
-- 3. Revertir cambios en memberships
-- =====================================================

-- 3.1 Eliminar metadata
ALTER TABLE memberships DROP COLUMN IF EXISTS metadata;

-- 3.2 Restaurar roles originales
ALTER TABLE memberships DROP CONSTRAINT IF EXISTS memberships_role_check;
ALTER TABLE memberships ADD CONSTRAINT memberships_role_check
    CHECK (role IN ('teacher', 'student', 'guardian'));

-- =====================================================
-- 4. Revertir cambios en schools
-- =====================================================

ALTER TABLE schools DROP COLUMN IF EXISTS metadata;

-- =====================================================
-- 5. Revertir cambios en academic_units
-- =====================================================

-- 5.1 Eliminar metadata y description
ALTER TABLE academic_units DROP COLUMN IF EXISTS metadata;
ALTER TABLE academic_units DROP COLUMN IF EXISTS description;

-- 5.2 Restaurar academic_year como NOT NULL
ALTER TABLE academic_units ALTER COLUMN academic_year SET NOT NULL;
ALTER TABLE academic_units ALTER COLUMN academic_year DROP DEFAULT;

-- 5.3 Restaurar tipos originales
ALTER TABLE academic_units DROP CONSTRAINT IF EXISTS academic_units_type_check;
ALTER TABLE academic_units ADD CONSTRAINT academic_units_type_check
    CHECK (type IN ('grade', 'class', 'section'));

-- 5.4 Eliminar constraint de auto-referencia
ALTER TABLE academic_units DROP CONSTRAINT IF EXISTS academic_units_no_self_reference;

-- 5.5 Eliminar índice de parent_unit_id
DROP INDEX IF EXISTS idx_academic_units_parent;

-- 5.6 Eliminar parent_unit_id
ALTER TABLE academic_units DROP COLUMN IF EXISTS parent_unit_id;

COMMIT;
