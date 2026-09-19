CREATE TABLE IF NOT EXISTS teacher_classes (
    teacher_id UUID NOT NULL REFERENCES users(id),
    class_id   UUID NOT NULL REFERENCES classes(id),
    PRIMARY KEY (teacher_id, class_id)
);
