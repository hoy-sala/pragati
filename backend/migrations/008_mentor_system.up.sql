CREATE TABLE mentor_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mentor_id UUID NOT NULL REFERENCES users(id),
    student_id UUID NOT NULL REFERENCES students(id),
    academic_year_id UUID NOT NULL REFERENCES academic_years(id),
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    assigned_by UUID REFERENCES users(id),
    UNIQUE(mentor_id, student_id, academic_year_id)
);

CREATE INDEX idx_mentor_assignments_mentor ON mentor_assignments(mentor_id);
CREATE INDEX idx_mentor_assignments_student ON mentor_assignments(student_id);
CREATE INDEX idx_mentor_assignments_year ON mentor_assignments(academic_year_id);

CREATE TABLE mentor_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mentor_id UUID NOT NULL REFERENCES users(id),
    student_id UUID NOT NULL REFERENCES students(id),
    date DATE NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'present' CHECK (status IN ('present', 'absent', 'late')),
    parent_contacted BOOLEAN DEFAULT false,
    parent_contact_time TIMESTAMPTZ,
    remarks TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(mentor_id, student_id, date)
);

CREATE INDEX idx_mentor_attendance_mentor_date ON mentor_attendance(mentor_id, date);
CREATE INDEX idx_mentor_attendance_student ON mentor_attendance(student_id);

CREATE TABLE mentor_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    mentor_id UUID NOT NULL REFERENCES users(id),
    student_id UUID NOT NULL REFERENCES students(id),
    log_date DATE NOT NULL DEFAULT CURRENT_DATE,
    category VARCHAR(20) NOT NULL CHECK (category IN ('health', 'behavior', 'grievance', 'academic')),
    severity VARCHAR(10) NOT NULL DEFAULT 'low' CHECK (severity IN ('low', 'medium', 'high', 'urgent')),
    description TEXT NOT NULL,
    action_taken TEXT,
    parent_informed BOOLEAN DEFAULT false,
    reviewed_by_principal BOOLEAN DEFAULT false,
    principal_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_mentor_logs_mentor ON mentor_logs(mentor_id);
CREATE INDEX idx_mentor_logs_student ON mentor_logs(student_id);
CREATE INDEX idx_mentor_logs_severity ON mentor_logs(severity) WHERE severity IN ('high', 'urgent');
CREATE INDEX idx_mentor_logs_date ON mentor_logs(log_date);
