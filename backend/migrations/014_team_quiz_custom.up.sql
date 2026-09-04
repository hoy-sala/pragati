ALTER TABLE team_quizzes ADD COLUMN IF NOT EXISTS team_names JSONB DEFAULT '[]';
-- timer_sec already exists, allow 0 for no timer
