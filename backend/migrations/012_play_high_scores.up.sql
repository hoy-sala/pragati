CREATE TABLE IF NOT EXISTS play_high_scores (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_name     VARCHAR(100) NOT NULL,
    class_id        UUID NOT NULL REFERENCES classes(id),
    subject_id      UUID NOT NULL REFERENCES subjects(id),
    topic           VARCHAR(200) NOT NULL DEFAULT '',
    difficulty      VARCHAR(20) NOT NULL DEFAULT 'medium',
    score           INTEGER NOT NULL,
    total_questions INTEGER NOT NULL,
    correct_count   INTEGER NOT NULL,
    best_streak     INTEGER NOT NULL DEFAULT 0,
    time_taken_ms   INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_play_scores_class_subject ON play_high_scores(class_id, subject_id);
CREATE INDEX idx_play_scores_difficulty ON play_high_scores(difficulty);
CREATE INDEX idx_play_scores_score ON play_high_scores(score DESC);
