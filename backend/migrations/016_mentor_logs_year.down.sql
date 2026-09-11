DROP INDEX IF EXISTS idx_mentor_logs_year;
ALTER TABLE mentor_logs DROP COLUMN IF EXISTS academic_year_id;
