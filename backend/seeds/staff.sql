-- Seed staff users
-- Default password: pragati123 (bcrypt hash pre-generated)
-- School ID: 00000000-0000-0000-0000-000000000001

INSERT INTO users (id, school_id, email, password_hash, name, role, phone, mobile, is_active)
VALUES
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9620078669@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Ajjaiah S', 'principal', '9620078669', '9620078669', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9986745650@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Siddappa N', 'teacher', '9986745650', '9986745650', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9845734969@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Manoj Kumar M J', 'teacher', '9845734969', '9845734969', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '8618536185@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Manoranjan E', 'teacher', '8618536185', '8618536185', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9742891290@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Vedamurthy D', 'teacher', '9742891290', '9742891290', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '8892920792@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Shivakumar K', 'teacher', '8892920792', '8892920792', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '8722255107@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Gopi K', 'teacher', '8722255107', '8722255107', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9986234946@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Kallesh K G', 'teacher', '9986234946', '9986234946', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9972404223@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Hoysala T', 'teacher', '9972404223', '9972404223', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9980999230@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Guruswamy E', 'teacher', '9980999230', '9980999230', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9071975205@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Shivayogi Kasavannavar', 'teacher', '9071975205', '9071975205', true),
  (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '9901960819@pragati.edu', '$2b$10$.IN8VNHrEQbH9zrCZctzs.NDiU7NJunmvjNq67okJt5CNPvE1fhhi', 'Sankranthi S Tilak', 'special_educator', '9901960819', '9901960819', true)
ON CONFLICT (email) DO NOTHING;
