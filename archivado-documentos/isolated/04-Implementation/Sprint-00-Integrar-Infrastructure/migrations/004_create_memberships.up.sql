-- Tabla: memberships (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-admin, api-mobile, worker
-- Versión: v0.7.0 (extendida con roles administrativos para api-admin)

CREATE TABLE IF NOT EXISTS memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    school_id UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    academic_unit_id UUID REFERENCES academic_units(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('teacher', 'student', 'guardian', 'coordinator', 'admin', 'assistant')),  -- EXTENDIDO
    metadata JSONB DEFAULT '{}'::jsonb,  -- NUEVO: Extensibilidad (permisos específicos, historial, etc.)
    is_active BOOLEAN NOT NULL DEFAULT true,
    enrolled_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    withdrawn_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, school_id, academic_unit_id, role)
);

CREATE INDEX idx_memberships_user ON memberships(user_id);
CREATE INDEX idx_memberships_school ON memberships(school_id);
CREATE INDEX idx_memberships_unit ON memberships(academic_unit_id);
CREATE INDEX idx_memberships_role ON memberships(role);
CREATE INDEX idx_memberships_active ON memberships(is_active);

COMMENT ON TABLE memberships IS 'Relación usuario-escuela-unidad académica';
COMMENT ON COLUMN memberships.role IS 'Rol: teacher, student, guardian, coordinator, admin, assistant';
COMMENT ON COLUMN memberships.metadata IS 'Metadata extensible: permisos específicos, configuración, historial';
COMMENT ON COLUMN memberships.enrolled_at IS 'Fecha de inicio de membresía';
COMMENT ON COLUMN memberships.withdrawn_at IS 'Fecha de fin de membresía (NULL = activo)';
