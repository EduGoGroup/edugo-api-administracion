-- Migration: 013_add_ltree_to_academic_units
-- Description: Add ltree extension and path column for efficient hierarchical queries
-- Dependencies: 012_extend_for_admin_api.up.sql
-- Date: 2025-11-18
-- Sprint: Sprint-03 - Repositorios con ltree

BEGIN;

-- =====================================================
-- 1. Enable ltree extension
-- =====================================================

CREATE EXTENSION IF NOT EXISTS ltree;

COMMENT ON EXTENSION ltree IS 'Data type for hierarchical tree-like structures';

-- =====================================================
-- 2. Add path column to academic_units
-- =====================================================

ALTER TABLE academic_units ADD COLUMN IF NOT EXISTS path ltree;

COMMENT ON COLUMN academic_units.path IS 'Materialized path using ltree for efficient hierarchical queries (e.g., ''unit_id1.unit_id2.unit_id3'')';

-- =====================================================
-- 3. Create indices for ltree performance
-- =====================================================

-- GIST index for ancestor/descendant queries (@>, <@, @)
CREATE INDEX IF NOT EXISTS academic_units_path_gist_idx ON academic_units USING GIST (path);

-- BTREE index for exact path lookups and sorting
CREATE INDEX IF NOT EXISTS academic_units_path_btree_idx ON academic_units USING btree (path);

COMMENT ON INDEX academic_units_path_gist_idx IS 'GIST index for ltree ancestor/descendant queries';
COMMENT ON INDEX academic_units_path_btree_idx IS 'BTREE index for exact path lookups and sorting';

-- =====================================================
-- 4. Create function to update path automatically
-- =====================================================

CREATE OR REPLACE FUNCTION update_academic_unit_path()
RETURNS TRIGGER AS $$
BEGIN
  -- If root unit (no parent), path is just the unit ID
  IF NEW.parent_unit_id IS NULL THEN
    NEW.path = NEW.id::text::ltree;
  ELSE
    -- If has parent, concatenate parent's path + unit ID
    SELECT path || NEW.id::text::ltree INTO NEW.path
    FROM academic_units
    WHERE id = NEW.parent_unit_id;

    -- Verify parent exists and has a path
    IF NEW.path IS NULL THEN
      RAISE EXCEPTION 'Parent unit % not found or has no path', NEW.parent_unit_id;
    END IF;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_academic_unit_path() IS 'Automatically updates path column when academic_unit is inserted or parent changes';

-- =====================================================
-- 5. Create trigger to call path update function
-- =====================================================

-- Drop old trigger if exists to avoid conflicts
DROP TRIGGER IF EXISTS academic_unit_path_trigger ON academic_units;

CREATE TRIGGER academic_unit_path_trigger
  BEFORE INSERT OR UPDATE OF parent_unit_id ON academic_units
  FOR EACH ROW
  EXECUTE FUNCTION update_academic_unit_path();

COMMENT ON TRIGGER academic_unit_path_trigger ON academic_units IS 'Trigger to maintain path column automatically';

-- =====================================================
-- 6. Update existing cycle prevention function to use ltree
-- =====================================================

-- Replace the old recursion-based cycle prevention with ltree-based version
CREATE OR REPLACE FUNCTION prevent_academic_unit_cycles()
RETURNS TRIGGER AS $$
BEGIN
  -- If no parent, no cycle possible
  IF NEW.parent_unit_id IS NULL THEN
    RETURN NEW;
  END IF;

  -- For UPDATE operations where unit already exists
  -- Check if the new parent would create a cycle
  IF TG_OP = 'UPDATE' AND NEW.id IS NOT NULL THEN
    -- Verify that the new parent is not a descendant of this unit
    -- Using ltree: if parent's path is contained in this unit's path, it's a cycle
    IF EXISTS (
      SELECT 1
      FROM academic_units
      WHERE id = NEW.parent_unit_id
        AND path <@ (SELECT path FROM academic_units WHERE id = NEW.id)
    ) THEN
      RAISE EXCEPTION 'Cannot set parent: would create a cycle in the hierarchy (parent % is descendant of unit %)',
        NEW.parent_unit_id, NEW.id;
    END IF;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION prevent_academic_unit_cycles() IS 'Prevents cycles in hierarchy using ltree path comparison (more efficient than recursion)';

-- =====================================================
-- 7. Populate path for existing records
-- =====================================================

-- For existing records, we need to calculate the path
-- This is done using a recursive CTE to build paths bottom-up
DO $$
DECLARE
  updated_count INTEGER;
BEGIN
  -- Only populate if there are records without paths
  IF EXISTS (SELECT 1 FROM academic_units WHERE path IS NULL LIMIT 1) THEN

    -- Update paths using recursive CTE
    WITH RECURSIVE unit_paths AS (
      -- Base case: root units (no parent)
      SELECT
        id,
        parent_unit_id,
        id::text::ltree AS computed_path
      FROM academic_units
      WHERE parent_unit_id IS NULL

      UNION ALL

      -- Recursive case: child units
      SELECT
        au.id,
        au.parent_unit_id,
        (up.computed_path || au.id::text)::ltree AS computed_path
      FROM academic_units au
      INNER JOIN unit_paths up ON au.parent_unit_id = up.id
    )
    UPDATE academic_units au
    SET path = up.computed_path
    FROM unit_paths up
    WHERE au.id = up.id;

    GET DIAGNOSTICS updated_count = ROW_COUNT;
    RAISE NOTICE 'Populated path for % existing academic_units', updated_count;
  ELSE
    RAISE NOTICE 'All academic_units already have paths, skipping population';
  END IF;
END $$;

-- =====================================================
-- 8. Make path NOT NULL after populating existing data
-- =====================================================

-- Now that all records have paths, make it required
ALTER TABLE academic_units ALTER COLUMN path SET NOT NULL;

COMMIT;
