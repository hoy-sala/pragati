CREATE TABLE IF NOT EXISTS assessment_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id   UUID NOT NULL REFERENCES schools(id),
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(20),
    weightage   NUMERIC(5,2) DEFAULT 0,
    sort_order  INT NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_assessment_categories_school ON assessment_categories(school_id) WHERE deleted_at IS NULL;

INSERT INTO assessment_categories (id, school_id, name, code, weightage, sort_order)
VALUES
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'FA1', 'FA1', 10, 1),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'FA2', 'FA2', 10, 2),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'FA3', 'FA3', 10, 3),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'FA4', 'FA4', 10, 4),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'SA1', 'SA1', 20, 5),
    (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'SA2', 'SA2', 20, 6);

CREATE TABLE IF NOT EXISTS assessments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    category_id     UUID NOT NULL REFERENCES assessment_categories(id),
    subject_id      UUID NOT NULL REFERENCES subjects(id),
    teacher_id      UUID REFERENCES users(id),
    class_id        UUID NOT NULL REFERENCES classes(id),
    section_id      UUID REFERENCES sections(id),
    name            VARCHAR(200),
    max_marks       NUMERIC(8,2) NOT NULL,
    weightage       NUMERIC(5,2) DEFAULT 100,
    date            DATE,
    chapters        JSONB DEFAULT '[]',
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    is_published    BOOLEAN NOT NULL DEFAULT FALSE,
    is_locked       BOOLEAN NOT NULL DEFAULT FALSE,
    version         INT NOT NULL DEFAULT 1,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_assessments_class ON assessments(class_id, section_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_assessments_subject ON assessments(subject_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_assessments_cat ON assessments(category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_assessments_year ON assessments(academic_year_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_assessments_teacher ON assessments(teacher_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS marks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id   UUID NOT NULL REFERENCES assessments(id),
    student_id      UUID NOT NULL REFERENCES students(id),
    marks_obtained  NUMERIC(8,2),
    is_absent       BOOLEAN NOT NULL DEFAULT FALSE,
    is_grace        BOOLEAN NOT NULL DEFAULT FALSE,
    remarks         TEXT,
    entered_by      UUID NOT NULL REFERENCES users(id),
    entered_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(assessment_id, student_id)
);

CREATE INDEX idx_marks_assessment ON marks(assessment_id);
CREATE INDEX idx_marks_student ON marks(student_id);
