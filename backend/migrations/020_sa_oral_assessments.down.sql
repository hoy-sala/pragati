DELETE FROM marks
WHERE assessment_id IN (
    SELECT a.id FROM assessments a
    JOIN classes c ON c.id = a.class_id
    JOIN subjects s ON s.id = a.subject_id
    WHERE a.name IN ('SA1 Oral', 'SA2 Oral')
      AND c.name IN ('Class 6', 'Class 7')
      AND s.code IN ('KAN', 'ENG', 'HIN', 'MAT', 'SCI', 'SOC')
);

DELETE FROM assessments
WHERE name IN ('SA1 Oral', 'SA2 Oral')
  AND deleted_at IS NULL
  AND class_id IN (SELECT id FROM classes WHERE name IN ('Class 6', 'Class 7'))
  AND subject_id IN (SELECT id FROM subjects WHERE code IN ('KAN', 'ENG', 'HIN', 'MAT', 'SCI', 'SOC'));
