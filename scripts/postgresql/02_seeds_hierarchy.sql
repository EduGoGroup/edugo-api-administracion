-- =====================================================
-- Script: 02_seeds_hierarchy.sql
-- Descripción: Seeds de datos de prueba para jerarquía académica
-- Proyecto: edugo-api-administracion
-- Fecha: 12 de Noviembre, 2025
-- =====================================================

-- Limpiar datos existentes (solo para testing)
TRUNCATE TABLE unit_membership CASCADE;
TRUNCATE TABLE academic_unit CASCADE;
TRUNCATE TABLE school CASCADE;

-- =====================================================
-- SEEDS: school
-- =====================================================
INSERT INTO school (id, name, code, address, contact_email, contact_phone) VALUES
    ('11111111-1111-1111-1111-111111111111', 'Colegio San José', 'ESC-001', 'Calle Principal 123, Ciudad', 'contacto@sanjose.edu', '+1-555-0001'),
    ('22222222-2222-2222-2222-222222222222', 'Instituto Nacional', 'ESC-002', 'Avenida Central 456, Ciudad', 'info@nacional.edu', '+1-555-0002'),
    ('33333333-3333-3333-3333-333333333333', 'Escuela Primaria Las Flores', 'ESC-003', 'Barrio Norte 789, Ciudad', 'admin@lasflores.edu', '+1-555-0003');

-- =====================================================
-- SEEDS: academic_unit
-- =====================================================

-- Colegio San José (ESC-001)
-- Nivel School (raíz)
INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('a0000000-0000-0000-0000-000000000001', NULL, '11111111-1111-1111-1111-111111111111', 'school', 'Colegio San José', 'SJ-ROOT');

-- Grados
INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('a1000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'grade', 'Primer Grado', 'SJ-G1'),
    ('a1000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'grade', 'Segundo Grado', 'SJ-G2'),
    ('a1000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'grade', 'Tercer Grado', 'SJ-G3');

-- Secciones de Primer Grado
INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('a2000000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'section', 'Primer Grado - Sección A', 'SJ-G1-A'),
    ('a2000000-0000-0000-0000-000000000002', 'a1000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'section', 'Primer Grado - Sección B', 'SJ-G1-B');

-- Secciones de Segundo Grado
INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('a2000000-0000-0000-0000-000000000003', 'a1000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'section', 'Segundo Grado - Sección A', 'SJ-G2-A'),
    ('a2000000-0000-0000-0000-000000000004', 'a1000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'section', 'Segundo Grado - Sección B', 'SJ-G2-B');

-- Clubs (paralelos a grados)
INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code, description) VALUES
    ('a3000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'club', 'Club de Robótica', 'SJ-CLUB-ROB', 'Club extracurricular de robótica'),
    ('a3000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'club', 'Club de Ajedrez', 'SJ-CLUB-AJE', 'Club de ajedrez para todos los niveles');

-- Departamentos administrativos
INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code, description) VALUES
    ('a4000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'department', 'Departamento de Matemáticas', 'SJ-DEPT-MAT', 'Coordinación académica de matemáticas'),
    ('a4000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'department', 'Departamento de Idiomas', 'SJ-DEPT-LANG', 'Coordinación académica de idiomas');

-- Instituto Nacional (ESC-002)
INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('b0000000-0000-0000-0000-000000000001', NULL, '22222222-2222-2222-2222-222222222222', 'school', 'Instituto Nacional', 'IN-ROOT');

INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('b1000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'grade', 'Primer Año', 'IN-Y1'),
    ('b1000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'grade', 'Segundo Año', 'IN-Y2');

INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('b2000000-0000-0000-0000-000000000001', 'b1000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', 'section', 'Primer Año - Sección A', 'IN-Y1-A');

-- Escuela Primaria Las Flores (ESC-003)
INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('c0000000-0000-0000-0000-000000000001', NULL, '33333333-3333-3333-3333-333333333333', 'school', 'Escuela Primaria Las Flores', 'LF-ROOT');

INSERT INTO academic_unit (id, parent_unit_id, school_id, unit_type, display_name, code) VALUES
    ('c1000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'grade', 'Preescolar', 'LF-PRE'),
    ('c1000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', 'grade', 'Primero', 'LF-G1');

-- =====================================================
-- SEEDS: unit_membership
-- =====================================================

-- UUIDs simulados de usuarios (en producción vendrían de tabla users)
-- Estudiantes en Primer Grado - Sección A (Colegio San José)
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
    ('a2000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'student', '2025-01-15'),
    ('a2000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'student', '2025-01-15'),
    ('a2000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000003', 'student', '2025-01-15');

-- Profesor en Primer Grado - Sección A
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
    ('a2000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000101', 'teacher', '2025-01-10');

-- Estudiantes en Primer Grado - Sección B
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
    ('a2000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000004', 'student', '2025-01-15'),
    ('a2000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000005', 'student', '2025-01-15');

-- Profesor en Primer Grado - Sección B
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
    ('a2000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000102', 'teacher', '2025-01-10');

-- Coordinador del Departamento de Matemáticas
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
    ('a4000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000201', 'coordinator', '2025-01-01');

-- Miembros del Club de Robótica
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
    ('a3000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'student', '2025-02-01'),
    ('a3000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000004', 'student', '2025-02-01'),
    ('a3000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000103', 'teacher', '2025-02-01');

-- Membresías con fecha de expiración (históricas)
INSERT INTO unit_membership (unit_id, user_id, role, valid_from, valid_until) VALUES
    ('a2000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000006', 'student', '2024-01-15', '2024-12-20');

-- Admin a nivel escuela (Colegio San José)
INSERT INTO unit_membership (unit_id, user_id, role, valid_from) VALUES
    ('a0000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000301', 'admin', '2024-08-01');

-- =====================================================
-- VERIFICACIÓN DE DATOS
-- =====================================================

-- Contar registros insertados
DO $$
DECLARE
    school_count INTEGER;
    unit_count INTEGER;
    membership_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO school_count FROM school;
    SELECT COUNT(*) INTO unit_count FROM academic_unit WHERE deleted_at IS NULL;
    SELECT COUNT(*) INTO membership_count FROM unit_membership;
    
    RAISE NOTICE 'Seeds insertados exitosamente:';
    RAISE NOTICE '  - Escuelas: %', school_count;
    RAISE NOTICE '  - Unidades académicas: %', unit_count;
    RAISE NOTICE '  - Membresías: %', membership_count;
END $$;
