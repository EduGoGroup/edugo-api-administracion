-- Tabla: schools (Owner: infrastructure)
-- Creada por: edugo-infrastructure
-- Usada por: api-admin, api-mobile, worker
-- Versión: v0.7.0 (extendida con metadata para api-admin)

CREATE TABLE IF NOT EXISTS schools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    address TEXT,
    city VARCHAR(100),
    country VARCHAR(100) NOT NULL DEFAULT 'Chile',
    phone VARCHAR(50),
    email VARCHAR(255),
    metadata JSONB DEFAULT '{}'::jsonb,  -- NUEVO: Extensibilidad (logo, config, etc.)
    is_active BOOLEAN NOT NULL DEFAULT true,
    subscription_tier VARCHAR(50) NOT NULL DEFAULT 'free' CHECK (subscription_tier IN ('free', 'basic', 'premium', 'enterprise')),
    max_teachers INTEGER NOT NULL DEFAULT 10,
    max_students INTEGER NOT NULL DEFAULT 100,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_schools_code ON schools(code);
CREATE INDEX idx_schools_active ON schools(is_active);
CREATE INDEX idx_schools_tier ON schools(subscription_tier);

COMMENT ON TABLE schools IS 'Escuelas/Instituciones educativas';
COMMENT ON COLUMN schools.subscription_tier IS 'Nivel de subscripción: free, basic, premium, enterprise';
COMMENT ON COLUMN schools.metadata IS 'Metadata extensible: logo, configuración institucional, etc.';
