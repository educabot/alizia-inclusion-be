package testutil

import (
	"time"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var TestOrgID = uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567890")

func Ptr[T any](v T) *T {
	return &v
}

func NewUser(id int64, name string) entities.User {
	return entities.User{
		ID:             id,
		OrganizationID: TestOrgID,
		Email:          name + "@test.com",
		Name:           name,
		Role:           "teacher",
		TimeTrackedEntity: entities.TimeTrackedEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func NewClassroom(id int64, name string) entities.Classroom {
	return entities.Classroom{
		ID:             id,
		OrganizationID: TestOrgID,
		Name:           name,
		TimeTrackedEntity: entities.TimeTrackedEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func NewStudent(id, classroomID int64, name string) entities.Student {
	return entities.Student{
		ID:             id,
		OrganizationID: TestOrgID,
		ClassroomID:    classroomID,
		Name:           name,
		TimeTrackedEntity: entities.TimeTrackedEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func NewStudentWithProfile(id, classroomID int64, name string, difficulties []string) entities.Student {
	s := NewStudent(id, classroomID, name)
	s.Profile = &entities.StudentProfile{
		ID:           id,
		StudentID:    id,
		IsTransitory: false,
		Difficulties: pq.StringArray(difficulties),
		TimeTrackedEntity: entities.TimeTrackedEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	return s
}

func NewRamp(id int64, name string) entities.Ramp {
	return entities.Ramp{
		ID:             id,
		OrganizationID: TestOrgID,
		Name:           name,
		SortOrder:      int(id),
		TimeTrackedEntity: entities.TimeTrackedEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func NewDevice(id, rampID int64, name string) entities.Device {
	return entities.Device{
		ID:             id,
		OrganizationID: TestOrgID,
		RampID:         rampID,
		Name:           name,
		Quantity:       1,
		SortOrder:      int(id),
		TimeTrackedEntity: entities.TimeTrackedEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func NewAdaptation(id, studentID, teacherID int64) entities.Adaptation {
	return entities.Adaptation{
		ID:             id,
		OrganizationID: TestOrgID,
		StudentID:      &studentID,
		TeacherID:      teacherID,
		Subject:        "Matematicas",
		Status:         "en_curso",
		TimeTrackedEntity: entities.TimeTrackedEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func NewConversation(id, userID int64, mode string) entities.Conversation {
	return entities.Conversation{
		ID:             id,
		OrganizationID: TestOrgID,
		UserID:         userID,
		Mode:           mode,
		TimeTrackedEntity: entities.TimeTrackedEntity{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
}

func NewAdaptationResource(id, adaptationID int64) entities.AdaptationResource {
	return entities.AdaptationResource{
		ID:           id,
		AdaptationID: adaptationID,
		Title:        "Resource",
		FileURL:      "https://example.com/file.pdf",
		FileType:     "pdf",
		CreatedAt:    time.Now(),
	}
}
