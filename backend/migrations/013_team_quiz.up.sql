CREATE TABLE IF NOT EXISTS team_quizzes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    title           VARCHAR(200) NOT NULL,
    description     TEXT,
    teams           INT NOT NULL CHECK (teams BETWEEN 2 AND 8),
    per_team        INT NOT NULL CHECK (per_team % 5 = 0 AND per_team BETWEEN 5 AND 20),
    timer_sec       INT NOT NULL DEFAULT 30,
    chapters        JSONB NOT NULL DEFAULT '[]',
    questions       JSONB NOT NULL DEFAULT '{}',
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_team_quizzes_school ON team_quizzes(school_id) WHERE deleted_at IS NULL;
