-- SA oral (10 marks) companion to SA1/SA2 for classes 6-7 curricular subjects.
-- Derived from the existing SA1/SA2 rows so school, category, subject, class,
-- teacher and academic year are inherited, and so the stray empty "Class 6"
-- row (which has no assessments) is never picked up as a source.
INSERT INTO assessments (
    id, school_id, category_id, subject_id, teacher_id, class_id,
    name, max_marks, academic_year_id, is_published, is_locked, version,
    created_at, updated_at
)
SELECT
    gen_random_uuid(), a.school_id, a.category_id, a.subject_id, a.teacher_id, a.class_id,
    v.oral_name, 10, a.academic_year_id, false, false, 1, NOW(), NOW()
FROM assessments a
JOIN classes c ON c.id = a.class_id
JOIN subjects s ON s.id = a.subject_id
CROSS JOIN (VALUES ('SA1', 'SA1 Oral'), ('SA2', 'SA2 Oral')) AS v(base_name, oral_name)
WHERE c.name IN ('Class 6', 'Class 7')
  AND s.code IN ('KAN', 'ENG', 'HIN', 'MAT', 'SCI', 'SOC')
  AND a.name = v.base_name
  AND a.deleted_at IS NULL
  AND NOT EXISTS (
      SELECT 1 FROM assessments o
      WHERE o.class_id = a.class_id
        AND o.subject_id = a.subject_id
        AND o.name = v.oral_name
        AND o.deleted_at IS NULL
  );
