-- Tabla: academic_units (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-admin (jerarquía), api-mobile (plano), worker
-- Versión: v0.7.0 (extendida con jerarquía para api-admin)

CREATE TABLE IF NOT EXISTS academic_units (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_unit_id UUID REFERENCES academic_units(id) ON DELETE SET NULL,  -- NUEVO: Jerarquía (NULL = raíz)
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('school', 'grade', 'class', 'section', 'club', 'department')),  -- EXTENDIDO
    description TEXT,  -- NUEVO: Descripción de la unidad
    level VARCHAR(50),
    academic_year INTEGER DEFAULT 0,  -- CAMBIADO: Ahora nullable, 0 = sin año específico
    metadata JSONB DEFAULT '{}'::jsonb,  -- NUEVO: Extensibilidad
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(school_id, code, academic_year),
    CONSTRAINT academic_units_no_self_reference CHECK (id != parent_unit_id)  -- NUEVO: Prevenir auto-referencia
);

CREATE INDEX idx_academic_units_parent ON academic_units(parent_unit_id);  -- NUEVO: Para jerarquía
CREATE INDEX idx_academic_units_school ON academic_units(school_id);
CREATE INDEX idx_academic_units_type ON academic_units(type);
CREATE INDEX idx_academic_units_year ON academic_units(academic_year);
CREATE INDEX idx_academic_units_active ON academic_units(is_active);

COMMENT ON TABLE academic_units IS 'Unidades académicas con soporte de jerarquía opcional';
COMMENT ON COLUMN academic_units.parent_unit_id IS 'Unidad padre (jerarquía: Facultad → Departamento). NULL = raíz';
COMMENT ON COLUMN academic_units.type IS 'Tipo: school, grade, class, section, club, department';
COMMENT ON COLUMN academic_units.description IS 'Descripción de la unidad académica';
COMMENT ON COLUMN academic_units.academic_year IS 'Año académico. 0 = sin año específico (usado por api-admin)';
COMMENT ON COLUMN academic_units.metadata IS 'Metadata extensible';

-- NUEVO: Función para prevenir ciclos en jerarquía
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

    -- Agregar el ID actual si existe
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

-- NUEVO: Trigger para prevenir ciclos
CREATE TRIGGER trigger_prevent_academic_unit_cycles
    BEFORE INSERT OR UPDATE OF parent_unit_id ON academic_units
    FOR EACH ROW
    EXECUTE FUNCTION prevent_academic_unit_cycles();

-- NUEVO: Vista para árbol jerárquico (CTE recursivo)
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
