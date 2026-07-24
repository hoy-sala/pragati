CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable trigram extension for fuzzy search on student names
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE IF NOT EXISTS schools (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(200) NOT NULL,
    code        VARCHAR(20) UNIQUE NOT NULL DEFAULT 'SCHOOL',
    logo_url    TEXT,
    config      JSONB DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

INSERT INTO schools (id, name, code) VALUES ('00000000-0000-0000-0000-000000000001', 'Default School', 'SCHOOL');

CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    name            VARCHAR(200) NOT NULL,
    role            VARCHAR(20) NOT NULL CHECK (role IN ('admin','principal','teacher','special_educator','student','parent')),
    phone           VARCHAR(20),
    avatar_url      TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_school ON users(school_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS academic_years (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id   UUID NOT NULL REFERENCES schools(id),
    name        VARCHAR(50) NOT NULL,
    start_date  DATE NOT NULL,
    end_date    DATE NOT NULL,
    is_current  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_academic_years_school ON academic_years(school_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS classes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    academic_year_id UUID REFERENCES academic_years(id),
    name            VARCHAR(50) NOT NULL,
    code            VARCHAR(20),
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_classes_school ON classes(school_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS sections (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id        UUID NOT NULL REFERENCES classes(id),
    name            VARCHAR(20) NOT NULL,
    class_teacher_id UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_sections_class ON sections(class_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS subjects (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id   UUID NOT NULL REFERENCES schools(id),
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(20),
    is_language BOOLEAN NOT NULL DEFAULT FALSE,
    is_core     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_subjects_school ON subjects(school_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS class_subjects (
    class_id    UUID NOT NULL REFERENCES classes(id),
    subject_id  UUID NOT NULL REFERENCES subjects(id),
    PRIMARY KEY (class_id, subject_id)
);

CREATE TABLE IF NOT EXISTS houses (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id   UUID NOT NULL REFERENCES schools(id),
    name        VARCHAR(50) NOT NULL,
    code        VARCHAR(10),
    color       VARCHAR(7),
    logo_url    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS students (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id         UUID NOT NULL REFERENCES schools(id),
    user_id           UUID REFERENCES users(id),
    sats_number       VARCHAR(9) NOT NULL,
    admission_no      VARCHAR(50),
    roll_no           INT,
    first_name        VARCHAR(100) NOT NULL,
    last_name         VARCHAR(100),
    date_of_birth     DATE,
    gender            VARCHAR(10),
    photo_url         TEXT,
    blood_group       VARCHAR(5),
    address           TEXT,
    phone             VARCHAR(20),
    email             VARCHAR(255),
    class_id          UUID NOT NULL REFERENCES classes(id),
    section_id        UUID REFERENCES sections(id),
    house_id          UUID REFERENCES houses(id),
    academic_year_id  UUID NOT NULL REFERENCES academic_years(id),
    parent_name       VARCHAR(200),
    parent_phone      VARCHAR(20),
    parent_email      VARCHAR(255),
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_students_sats ON students(sats_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_class ON students(class_id, section_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_year ON students(academic_year_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_students_name ON students USING gin (first_name gin_trgm_ops, last_name gin_trgm_ops);

CREATE TABLE IF NOT EXISTS teacher_subjects (
    teacher_id  UUID NOT NULL REFERENCES users(id),
    subject_id  UUID NOT NULL REFERENCES subjects(id),
    class_id    UUID NOT NULL REFERENCES classes(id),
    PRIMARY KEY (teacher_id, subject_id, class_id)
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id),
    token_hash  VARCHAR(255) NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at  TIMESTAMPTZ
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
