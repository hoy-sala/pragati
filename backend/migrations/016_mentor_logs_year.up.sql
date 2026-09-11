ALTER TABLE mentor_logs ADD COLUMN IF NOT EXISTS academic_year_id UUID REFERENCES academic_years(id);
UPDATE mentor_logs ml SET academic_year_id = s.academic_year_id
FROM students s WHERE s.id = ml.student_id AND ml.academic_year_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_mentor_logs_year ON mentor_logs(academic_year_id);
