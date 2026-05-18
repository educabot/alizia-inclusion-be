CREATE TABLE adaptations (
    id BIGSERIAL PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    student_id BIGINT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    teacher_id BIGINT NOT NULL REFERENCES users(id),
    device_id BIGINT REFERENCES devices(id),
    subject VARCHAR(255) NOT NULL,
    activity_description TEXT,
    adaptation_strategy TEXT,
    outcome TEXT,
    notes TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_adaptations_organization_id ON adaptations(organization_id);
CREATE INDEX idx_adaptations_student_id ON adaptations(student_id);
CREATE INDEX idx_adaptations_teacher_id ON adaptations(teacher_id);
