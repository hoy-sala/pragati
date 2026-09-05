CREATE TABLE IF NOT EXISTS map_places (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id       UUID NOT NULL REFERENCES schools(id),
    name            VARCHAR(200) NOT NULL,
    kind            VARCHAR(50) NOT NULL DEFAULT 'place',
    category        VARCHAR(100) NOT NULL,
    map             VARCHAR(20) NOT NULL DEFAULT 'india' CHECK (map IN ('india','karnataka')),
    lat             NUMERIC(9,5),
    lng             NUMERIC(9,5),
    svg_x           NUMERIC(9,2),
    svg_y           NUMERIC(9,2),
    state           VARCHAR(100),
    district        VARCHAR(100),
    why_in_news     TEXT,
    news_date       DATE,
    exam_tags       JSONB NOT NULL DEFAULT '[]',
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT uq_map_places_name_map UNIQUE (name, map)
);

CREATE INDEX IF NOT EXISTS idx_map_places_school ON map_places(school_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_map_places_category ON map_places(category) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_map_places_map ON map_places(map) WHERE deleted_at IS NULL;
