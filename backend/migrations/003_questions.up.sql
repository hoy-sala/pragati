CREATE TABLE IF NOT EXISTS questions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    subject_id      UUID NOT NULL REFERENCES subjects(id),
    teacher_id      UUID REFERENCES users(id),
    question_type   VARCHAR(20) NOT NULL CHECK (question_type IN ('mcq','true_false','fill_blank','short_answer','long_answer')),
    question_text   TEXT NOT NULL,
    question_image  TEXT,
    options         JSONB DEFAULT '[]',
    answer          TEXT NOT NULL,
    marks           NUMERIC(6,2) NOT NULL DEFAULT 1,
    difficulty      VARCHAR(20) NOT NULL DEFAULT 'medium' CHECK (difficulty IN ('easy','medium','hard')),
    chapters        JSONB DEFAULT '[]',
    tags            JSONB DEFAULT '[]',
    explanation     TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_questions_subject ON questions(subject_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_questions_type ON questions(question_type) WHERE deleted_at IS NULL;
CREATE INDEX idx_questions_difficulty ON questions(difficulty) WHERE deleted_at IS NULL;
CREATE INDEX idx_questions_school ON questions(school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_questions_tags ON questions USING gin(tags);
