-- Holistic Progress Card (HPC) schema
-- Aligned with NEP 2020 PARAKH framework & Karnataka DSERT format

CREATE TABLE IF NOT EXISTS hpc_configurations (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id        UUID NOT NULL REFERENCES schools(id),
    stage            VARCHAR(20) NOT NULL CHECK (stage IN ('foundational','preparatory','middle','secondary')),
    class_id         UUID REFERENCES classes(id),
    academic_year_id UUID REFERENCES academic_years(id),
    grading_scheme   JSONB NOT NULL DEFAULT '[
        {"grade":"A+","min_pct":90,"max_pct":100,"descriptor":"Outstanding"},
        {"grade":"A","min_pct":75,"max_pct":89,"descriptor":"Excellent"},
        {"grade":"B+","min_pct":60,"max_pct":74,"descriptor":"Very Good"},
        {"grade":"B","min_pct":45,"max_pct":59,"descriptor":"Good"},
        {"grade":"C","min_pct":33,"max_pct":44,"descriptor":"Satisfactory"},
        {"grade":"D","min_pct":20,"max_pct":32,"descriptor":"Needs Improvement"},
        {"grade":"E","min_pct":0,"max_pct":19,"descriptor":"Requires Remedial"}
    ]',
    proficiency_scale JSONB NOT NULL DEFAULT '[
        {"level":1,"label":"Beginning","descriptor":"Needs significant support"},
        {"level":2,"label":"Developing","descriptor":"Occasional support needed"},
        {"level":3,"label":"Proficient","descriptor":"Meets expectations independently"},
        {"level":4,"label":"Advanced","descriptor":"Exceeds expectations, can guide peers"}
    ]',
    co_scholastic_areas JSONB NOT NULL DEFAULT '{
        "life_skills": {
            "thinking": {"label":"Thinking Skills","max_grade":"A"},
            "social": {"label":"Social Skills","max_grade":"A"},
            "emotional": {"label":"Emotional Skills","max_grade":"A"}
        },
        "attitudes": {
            "towards_teachers": {"label":"Towards Teachers","max_grade":"A"},
            "towards_schoolmates": {"label":"Towards Schoolmates","max_grade":"A"},
            "towards_school": {"label":"Towards School Programs","max_grade":"A"},
            "towards_environment": {"label":"Towards Environment","max_grade":"A"}
        },
        "values": {
            "honesty": {"label":"Honesty","max_grade":"A"},
            "responsibility": {"label":"Responsibility","max_grade":"A"},
            "cooperation": {"label":"Cooperation","max_grade":"A"},
            "discipline": {"label":"Discipline","max_grade":"A"}
        },
        "participation": {
            "sports": {"label":"Sports & Games","max_grade":"A"},
            "arts": {"label":"Arts & Crafts","max_grade":"A"},
            "music": {"label":"Music & Dance","max_grade":"A"},
            "clubs": {"label":"Clubs & Societies","max_grade":"A"}
        }
    }',
    health_params JSONB NOT NULL DEFAULT '{
        "track_height": true,
        "track_weight": true,
        "track_bmi": true,
        "track_vision": true,
        "track_dental": true,
        "track_physical_fitness": true
    }',
    terms JSONB NOT NULL DEFAULT '["Term1","Term2"]',
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_hpc_config_unique
    ON hpc_configurations(school_id, COALESCE(class_id, '00000000-0000-0000-0000-000000000000'))
    WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS learning_outcomes (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id        UUID NOT NULL REFERENCES schools(id),
    subject_id       UUID NOT NULL REFERENCES subjects(id),
    class_id         UUID NOT NULL REFERENCES classes(id),
    code             VARCHAR(20) NOT NULL,
    description      TEXT NOT NULL,
    domain           VARCHAR(20) NOT NULL CHECK (domain IN ('cognitive','affective','psychomotor')),
    expected_level   INT NOT NULL DEFAULT 3 CHECK (expected_level BETWEEN 1 AND 4),
    sort_order       INT NOT NULL DEFAULT 0,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_lo_code_unique
    ON learning_outcomes(school_id, subject_id, class_id, code)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_lo_subject_class ON learning_outcomes(subject_id, class_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS hpc_entries (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id        UUID NOT NULL REFERENCES students(id),
    academic_year_id  UUID NOT NULL REFERENCES academic_years(id),
    term              VARCHAR(20) NOT NULL,
    status            VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft','published')),
    scholastic        JSONB NOT NULL DEFAULT '[]',
    co_scholastic     JSONB NOT NULL DEFAULT '{}',
    health_pe         JSONB NOT NULL DEFAULT '{}',
    work_education    JSONB NOT NULL DEFAULT '{}',
    self_assessment   JSONB NOT NULL DEFAULT '{}',
    peer_assessment   JSONB NOT NULL DEFAULT '{}',
    parent_feedback   JSONB NOT NULL DEFAULT '{}',
    teacher_remarks   TEXT NOT NULL DEFAULT '',
    attendance_summary JSONB NOT NULL DEFAULT '{}',
    generated_pdf_url TEXT NOT NULL DEFAULT '',
    version           INT NOT NULL DEFAULT 1,
    locked_at         TIMESTAMPTZ,
    locked_by         UUID REFERENCES users(id),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_hpc_entry_unique
    ON hpc_entries(student_id, academic_year_id, term)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_hpc_entry_student ON hpc_entries(student_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_hpc_entry_term ON hpc_entries(academic_year_id, term) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS learning_outcome_assessments (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id         UUID NOT NULL REFERENCES students(id),
    learning_outcome_id UUID NOT NULL REFERENCES learning_outcomes(id),
    subject_id         UUID NOT NULL REFERENCES subjects(id),
    term               VARCHAR(20) NOT NULL,
    proficiency_level  INT NOT NULL CHECK (proficiency_level BETWEEN 1 AND 4),
    assessed_by        UUID NOT NULL REFERENCES users(id),
    assessment_date    DATE NOT NULL DEFAULT CURRENT_DATE,
    remarks            TEXT NOT NULL DEFAULT '',
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at         TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_loa_unique
    ON learning_outcome_assessments(student_id, learning_outcome_id, term)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_loa_student_term ON learning_outcome_assessments(student_id, term) WHERE deleted_at IS NULL;
CREATE INDEX idx_loa_subject_term ON learning_outcome_assessments(subject_id, term) WHERE deleted_at IS NULL;
