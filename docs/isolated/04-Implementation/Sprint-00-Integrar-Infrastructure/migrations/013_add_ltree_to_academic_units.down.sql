-- Migration Down: 013_add_ltree_to_academic_units
-- Description: Remove ltree extension and path column from academic_units
-- Date: 2025-11-18
-- Sprint: Sprint-03 - Repositorios con ltree

BEGIN;

-- =====================================================
-- 1. Drop triggers
-- =====================================================

DROP TRIGGER IF EXISTS academic_unit_path_trigger ON academic_units;

-- =====================================================
-- 2. Restore original cycle prevention function
-- =====================================================

-- Restore the recursion-based cycle prevention (original version from 012)
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

-- =====================================================
-- 3. Drop ltree-specific function
-- =====================================================

DROP FUNCTION IF EXISTS update_academic_unit_path();

-- =====================================================
-- 4. Drop indices
-- =====================================================

DROP INDEX IF EXISTS academic_units_path_btree_idx;
DROP INDEX IF EXISTS academic_units_path_gist_idx;

-- =====================================================
-- 5. Drop path column
-- =====================================================

ALTER TABLE academic_units DROP COLUMN IF EXISTS path;

-- =====================================================
-- 6. Drop ltree extension (if no other tables use it)
-- =====================================================

-- Note: Only drop if no other tables are using ltree
-- Uncomment the following line if you're sure no other tables use ltree:
-- DROP EXTENSION IF EXISTS ltree;

COMMENT ON EXTENSION ltree IS 'Extension preserved - remove manually if not used elsewhere';

COMMIT;
