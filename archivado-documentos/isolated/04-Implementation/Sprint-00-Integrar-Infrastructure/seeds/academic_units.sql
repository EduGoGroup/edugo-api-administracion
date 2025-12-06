-- Seeds de academic_units con jerarquía para api-admin
-- Ejecutar después de seeds de schools
-- Demuestra estructura jerárquica: Escuela → Facultad → Departamento → Carrera

-- ==============================================================
-- Escuela 1: Liceo Técnico Santiago (LTS-001)
-- Estructura: Escuela → Grados → Secciones
-- ==============================================================

-- Raíz: La escuela misma
INSERT INTO academic_units (id, parent_unit_id, school_id, name, code, type, academic_year) VALUES
('a1000000-0000-0000-0000-000000000001', NULL, '44444444-4444-4444-4444-444444444444', 'Liceo Técnico Santiago', 'LTS-ROOT', 'school', 0)
ON CONFLICT DO NOTHING;

-- Nivel 1: Grados
INSERT INTO academic_units (id, parent_unit_id, school_id, name, code, type, description, academic_year) VALUES
('a1100000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', '1° Medio', 'LTS-1M', 'grade', 'Primer año de educación media', 2025),
('a1200000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', '2° Medio', 'LTS-2M', 'grade', 'Segundo año de educación media', 2025),
('a1300000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', '3° Medio', 'LTS-3M', 'grade', 'Tercer año de educación media', 2025)
ON CONFLICT DO NOTHING;

-- Nivel 2: Secciones (clases) de 1° Medio
INSERT INTO academic_units (id, parent_unit_id, school_id, name, code, type, description, academic_year) VALUES
('a1110000-0000-0000-0000-000000000001', 'a1100000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', '1° Medio A', 'LTS-1M-A', 'section', 'Sección A de 1° Medio', 2025),
('a1120000-0000-0000-0000-000000000001', 'a1100000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', '1° Medio B', 'LTS-1M-B', 'section', 'Sección B de 1° Medio', 2025)
ON CONFLICT DO NOTHING;

-- Nivel 2: Clubs/Actividades extracurriculares
INSERT INTO academic_units (id, parent_unit_id, school_id, name, code, type, description, academic_year, metadata) VALUES
('a1400000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', 'Club de Robótica', 'LTS-ROBOT', 'club', 'Club de robótica educativa', 0, '{"schedule": "Martes 15:00-17:00", "max_students": 25}'),
('a1500000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001', '44444444-4444-4444-4444-444444444444', 'Club de Ciencias', 'LTS-CIENC', 'club', 'Club de ciencias experimentales', 0, '{"schedule": "Jueves 15:00-17:00", "max_students": 30}')
ON CONFLICT DO NOTHING;

-- ==============================================================
-- Escuela 2: Colegio Valparaíso (CV-002)
-- Estructura: Escuela → Departamentos → Clases
-- ==============================================================

-- Raíz: La escuela misma
INSERT INTO academic_units (id, parent_unit_id, school_id, name, code, type, academic_year) VALUES
('a2000000-0000-0000-0000-000000000002', NULL, '55555555-5555-5555-5555-555555555555', 'Colegio Valparaíso', 'CV-ROOT', 'school', 0)
ON CONFLICT DO NOTHING;

-- Nivel 1: Departamentos académicos
INSERT INTO academic_units (id, parent_unit_id, school_id, name, code, type, description, academic_year, metadata) VALUES
('a2100000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002', '55555555-5555-5555-5555-555555555555', 'Departamento de Matemáticas', 'CV-DMAT', 'department', 'Departamento de ciencias matemáticas', 0, '{"head_teacher": "Prof. González"}'),
('a2200000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002', '55555555-5555-5555-5555-555555555555', 'Departamento de Lenguaje', 'CV-DLEN', 'department', 'Departamento de lengua y literatura', 0, '{"head_teacher": "Prof. Martínez"}')
ON CONFLICT DO NOTHING;

-- Nivel 2: Clases del departamento de Matemáticas
INSERT INTO academic_units (id, parent_unit_id, school_id, name, code, type, description, academic_year) VALUES
('a2110000-0000-0000-0000-000000000002', 'a2100000-0000-0000-0000-000000000002', '55555555-5555-5555-5555-555555555555', 'Álgebra I', 'CV-MAT-ALG1', 'class', 'Álgebra nivel básico', 2025),
('a2120000-0000-0000-0000-000000000002', 'a2100000-0000-0000-0000-000000000002', '55555555-5555-5555-5555-555555555555', 'Geometría', 'CV-MAT-GEO', 'class', 'Geometría euclidiana', 2025)
ON CONFLICT DO NOTHING;

-- ==============================================================
-- Validar jerarquía (queries de prueba)
-- ==============================================================

-- Consulta 1: Ver toda la jerarquía usando la vista
-- SELECT * FROM v_academic_unit_tree ORDER BY school_id, path;

-- Consulta 2: Ver hijos de una unidad específica
-- SELECT * FROM academic_units WHERE parent_unit_id = 'a1000000-0000-0000-0000-000000000001';

-- Consulta 3: Ver path completo de una unidad
-- WITH RECURSIVE ancestors AS (
--     SELECT id, parent_unit_id, name, 1 as level
--     FROM academic_units
--     WHERE id = 'a1110000-0000-0000-0000-000000000001'
--     UNION ALL
--     SELECT au.id, au.parent_unit_id, au.name, a.level + 1
--     FROM academic_units au
--     INNER JOIN ancestors a ON au.id = a.parent_unit_id
-- )
-- SELECT * FROM ancestors ORDER BY level DESC;
