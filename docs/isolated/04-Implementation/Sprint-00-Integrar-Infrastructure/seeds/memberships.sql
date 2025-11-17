-- Seeds de memberships para api-admin
-- Ejecutar después de seeds de users y academic_units
-- Demuestra diferentes roles: teacher, student, coordinator, admin

-- ==============================================================
-- Memberships para Liceo Técnico Santiago
-- ==============================================================

-- Admin de prueba (11111111-1111-1111-1111-111111111111) como ADMIN de la escuela
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active) VALUES
('m1000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', '44444444-4444-4444-4444-444444444444', 'a1000000-0000-0000-0000-000000000001', 'admin', true)
ON CONFLICT DO NOTHING;

-- Docente de prueba (22222222-2222-2222-2222-222222222222) como TEACHER en 1° Medio A
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, metadata) VALUES
('m2000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', '44444444-4444-4444-4444-444444444444', 'a1110000-0000-0000-0000-000000000001', 'teacher', true, '{"subject": "Matemáticas", "hours_per_week": 4}')
ON CONFLICT DO NOTHING;

-- Mismo docente como COORDINATOR del Club de Robótica
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, metadata) VALUES
('m2100000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', '44444444-4444-4444-4444-444444444444', 'a1400000-0000-0000-0000-000000000001', 'coordinator', true, '{"position": "Coordinador de Club"}')
ON CONFLICT DO NOTHING;

-- Estudiante de prueba (33333333-3333-3333-3333-333333333333) como STUDENT en 1° Medio A
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active) VALUES
('m3000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444', 'a1110000-0000-0000-0000-000000000001', 'student', true)
ON CONFLICT DO NOTHING;

-- Mismo estudiante en Club de Robótica
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active) VALUES
('m3100000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444', 'a1400000-0000-0000-0000-000000000001', 'student', true)
ON CONFLICT DO NOTHING;

-- ==============================================================
-- Memberships para Colegio Valparaíso
-- ==============================================================

-- Admin como ASSISTANT (ayudante administrativo)
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, metadata) VALUES
('m4000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', '55555555-5555-5555-5555-555555555555', 'a2000000-0000-0000-0000-000000000002', 'assistant', true, '{"position": "Secretaría Académica"}')
ON CONFLICT DO NOTHING;

-- ==============================================================
-- Memberships inactivas (ejemplos de withdrawn)
-- ==============================================================

-- Estudiante que se retiró del club
INSERT INTO memberships (id, user_id, school_id, academic_unit_id, role, is_active, enrolled_at, withdrawn_at) VALUES
('m5000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333333', '44444444-4444-4444-4444-444444444444', 'a1500000-0000-0000-0000-000000000001', 'student', false, '2025-03-01 00:00:00+00', '2025-06-15 00:00:00+00')
ON CONFLICT DO NOTHING;

-- ==============================================================
-- Queries de validación
-- ==============================================================

-- Consulta 1: Ver memberships activas por escuela
-- SELECT
--     u.first_name, u.last_name,
--     s.name as school_name,
--     au.name as unit_name,
--     m.role
-- FROM memberships m
-- JOIN users u ON m.user_id = u.id
-- JOIN schools s ON m.school_id = s.id
-- LEFT JOIN academic_units au ON m.academic_unit_id = au.id
-- WHERE m.is_active = true
-- ORDER BY s.name, m.role, u.last_name;

-- Consulta 2: Ver todos los roles de un usuario
-- SELECT
--     s.name as school,
--     au.name as unit,
--     m.role,
--     m.is_active
-- FROM memberships m
-- JOIN schools s ON m.school_id = s.id
-- LEFT JOIN academic_units au ON m.academic_unit_id = au.id
-- WHERE m.user_id = '22222222-2222-2222-2222-222222222222'
-- ORDER BY s.name, au.name;

-- Consulta 3: Ver estudiantes de una unidad
-- SELECT
--     u.first_name, u.last_name, u.email,
--     m.enrolled_at
-- FROM memberships m
-- JOIN users u ON m.user_id = u.id
-- WHERE m.academic_unit_id = 'a1110000-0000-0000-0000-000000000001'
--   AND m.role = 'student'
--   AND m.is_active = true;
