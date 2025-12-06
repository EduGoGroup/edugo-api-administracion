DROP VIEW IF EXISTS v_academic_unit_tree;
DROP TRIGGER IF EXISTS trigger_prevent_academic_unit_cycles ON academic_units;
DROP FUNCTION IF EXISTS prevent_academic_unit_cycles();
DROP INDEX IF EXISTS idx_academic_units_active;
DROP INDEX IF EXISTS idx_academic_units_year;
DROP INDEX IF EXISTS idx_academic_units_type;
DROP INDEX IF EXISTS idx_academic_units_school;
DROP INDEX IF EXISTS idx_academic_units_parent;
DROP TABLE IF EXISTS academic_units;
