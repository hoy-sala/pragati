CREATE TABLE IF NOT EXISTS certificate_events (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id         UUID NOT NULL REFERENCES schools(id),
    academic_year_id   UUID REFERENCES academic_years(id),
    name              VARCHAR(200) NOT NULL,
    category          VARCHAR(50) NOT NULL DEFAULT 'other',
    held_date         DATE,
    venue             VARCHAR(200),
    description       TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_certificate_events_school ON certificate_events(school_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS certificate_signatories (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id      UUID NOT NULL REFERENCES certificate_events(id),
    name          VARCHAR(200) NOT NULL,
    role          VARCHAR(50) NOT NULL DEFAULT 'judge',
    title         VARCHAR(200),
    signature_url TEXT,
    sort_order    INT NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_certificate_signatories_event ON certificate_signatories(event_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS certificates (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id     UUID NOT NULL REFERENCES schools(id),
    event_id      UUID NOT NULL REFERENCES certificate_events(id),
    student_id    UUID NOT NULL REFERENCES students(id),
    position      VARCHAR(20) NOT NULL DEFAULT 'participation',
    prize_title   VARCHAR(100),
    issue_date    DATE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT uq_certificates_event_student UNIQUE (event_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_certificates_school ON certificates(school_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_certificates_event ON certificates(event_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_certificates_student ON certificates(student_id) WHERE deleted_at IS NULL;