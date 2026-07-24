-- Quiz/Online Test schema
-- Supports both student academic testing and staff training/assessment

CREATE TABLE IF NOT EXISTS quiz_assignments (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id        UUID NOT NULL REFERENCES schools(id),
    title            VARCHAR(200) NOT NULL,
    description      TEXT NOT NULL DEFAULT '',
    target_type      VARCHAR(20) NOT NULL CHECK (target_type IN ('student', 'staff')),
    target_id        UUID,
    pass_pct         INT NOT NULL DEFAULT 40 CHECK (pass_pct BETWEEN 0 AND 100),
    max_attempts     INT NOT NULL DEFAULT 1 CHECK (max_attempts > 0),
    duration_min     INT CHECK (duration_min > 0),
    shuffle_questions BOOLEAN NOT NULL DEFAULT TRUE,
    shuffle_options  BOOLEAN NOT NULL DEFAULT TRUE,
    show_result      BOOLEAN NOT NULL DEFAULT TRUE,
    start_at         TIMESTAMPTZ,
    end_at           TIMESTAMPTZ,
    is_published     BOOLEAN NOT NULL DEFAULT FALSE,
    created_by       UUID NOT NULL REFERENCES users(id),
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_quiz_school ON quiz_assignments(school_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_quiz_target ON quiz_assignments(target_type, target_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS quiz_questions (
    quiz_id      UUID NOT NULL REFERENCES quiz_assignments(id) ON DELETE CASCADE,
    question_id  UUID NOT NULL REFERENCES questions(id),
    sort_order   INT NOT NULL DEFAULT 0,
    marks        INT NOT NULL DEFAULT 1 CHECK (marks > 0),
    PRIMARY KEY (quiz_id, question_id)
);

CREATE INDEX idx_qq_quiz ON quiz_questions(quiz_id);

CREATE TABLE IF NOT EXISTS quiz_attempts (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id          UUID NOT NULL REFERENCES quiz_assignments(id),
    user_id          UUID NOT NULL REFERENCES users(id),
    attempt_no       INT NOT NULL DEFAULT 1,
    status           VARCHAR(20) NOT NULL DEFAULT 'in_progress'
                     CHECK (status IN ('in_progress','submitted','timed_out','graded')),
    score            NUMERIC(5,2),
    percentage       NUMERIC(5,2),
    passed           BOOLEAN,
    started_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    submitted_at     TIMESTAMPTZ,
    graded_at        TIMESTAMPTZ,
    graded_by        UUID REFERENCES users(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_attempt_unique
    ON quiz_attempts(quiz_id, user_id, attempt_no);

CREATE INDEX idx_attempt_user ON quiz_attempts(user_id);
CREATE INDEX idx_attempt_quiz ON quiz_attempts(quiz_id);

CREATE TABLE IF NOT EXISTS quiz_responses (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id       UUID NOT NULL REFERENCES quiz_attempts(id) ON DELETE CASCADE,
    question_id      UUID NOT NULL REFERENCES questions(id),
    selected_options TEXT[],
    text_answer      TEXT NOT NULL DEFAULT '',
    is_correct       BOOLEAN,
    marks_awarded    NUMERIC(5,2) NOT NULL DEFAULT 0,
    marks_total      NUMERIC(5,2) NOT NULL DEFAULT 1,
    graded_at        TIMESTAMPTZ,
    graded_by        UUID REFERENCES users(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (attempt_id, question_id)
);

CREATE INDEX idx_qr_attempt ON quiz_responses(attempt_id);
