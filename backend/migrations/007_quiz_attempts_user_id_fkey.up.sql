-- quiz_attempts.user_id can reference either users (staff) or students (student logins)
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS quiz_attempts_user_id_fkey;
