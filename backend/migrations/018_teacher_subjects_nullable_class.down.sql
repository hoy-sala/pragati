ALTER TABLE teacher_subjects DROP CONSTRAINT IF EXISTS teacher_subjects_pkey;
ALTER TABLE teacher_subjects ALTER COLUMN class_id SET NOT NULL;
ALTER TABLE teacher_subjects ADD CONSTRAINT teacher_subjects_pkey PRIMARY KEY (teacher_id, subject_id, class_id);
