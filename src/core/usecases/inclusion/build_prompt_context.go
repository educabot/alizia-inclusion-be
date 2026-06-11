package inclusion

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// maxPastAdaptations caps how many past adaptations are included in the prompt context.
const maxPastAdaptations = 10

// Missing-data labels — what Alizia may suggest completing, never require.
const (
	missingTeacherProfile = "perfil_docente"
	missingStudentProfile = "perfil_alumno"
	missingPPI            = "ppi"
	missingDiagnoses      = "diagnosticos"
)

// PromptContext is the typed context fed into the prompt. Ordered for caching:
// static fields first (eager, org-level cacheable), dynamic fields second (lazy,
// driven by the Prompt 0 router dimension). All dynamic blocks are optional —
// Alizia operates with whatever is available.
type PromptContext struct {
	// ---- STATIC (eager, cacheable) ----
	DeviceCatalog []entities.Device    `json:"device_catalog"`
	Situations    []entities.Situation `json:"situations"`

	// ---- DYNAMIC (lazy, driven by dimension) ----
	Dimension         string                         `json:"dimension"`
	Teacher           *entities.TeacherProfile       `json:"teacher,omitempty"`
	Classroom         *entities.Classroom            `json:"classroom,omitempty"`
	ClassroomStudents []entities.Student             `json:"classroom_students,omitempty"`
	TargetStudent     *entities.Student              `json:"target_student,omitempty"`
	Diagnoses         []entities.StudentDiagnosis    `json:"diagnoses,omitempty"`
	PPI               *entities.PPI                  `json:"ppi,omitempty"`
	PastAdaptations   []entities.Adaptation          `json:"past_adaptations,omitempty"`
	PriorSummaries    []entities.ConversationSummary `json:"prior_summaries,omitempty"`

	// MissingData lists by code which optional data is absent, so Alizia can suggest
	// completing it. Never blocks execution or surfaces as "N/A" in output.
	MissingData []string `json:"missing_data,omitempty"`
}

// ContextSnapshot holds a PII-free snapshot of the context (IDs only) for
// tracing in ai_usage / logs without exposing names or diagnoses.
type ContextSnapshot struct {
	Dimension          string   `json:"dimension"`
	TeacherUserID      *int64   `json:"teacher_user_id,omitempty"`
	ClassroomID        *int64   `json:"classroom_id,omitempty"`
	TargetStudentID    *int64   `json:"target_student_id,omitempty"`
	StudentProfileID   *int64   `json:"student_profile_id,omitempty"`
	PPIID              *int64   `json:"ppi_id,omitempty"`
	DiagnosesCount     int      `json:"diagnoses_count"`
	PastAdaptationsLen int      `json:"past_adaptations_len"`
	PriorSummariesLen  int      `json:"prior_summaries_len"`
	MissingData        []string `json:"missing_data,omitempty"`
}

// Snapshot builds a PII-free fingerprint of the context for traceability.
func (c *PromptContext) Snapshot() ContextSnapshot {
	snap := ContextSnapshot{
		Dimension:          c.Dimension,
		DiagnosesCount:     len(c.Diagnoses),
		PastAdaptationsLen: len(c.PastAdaptations),
		PriorSummariesLen:  len(c.PriorSummaries),
		MissingData:        c.MissingData,
	}
	if c.Teacher != nil {
		snap.TeacherUserID = &c.Teacher.UserID
	}
	if c.Classroom != nil {
		snap.ClassroomID = &c.Classroom.ID
	}
	if c.TargetStudent != nil {
		snap.TargetStudentID = &c.TargetStudent.ID
		if c.TargetStudent.Profile != nil {
			snap.StudentProfileID = &c.TargetStudent.Profile.ID
		}
	}
	if c.PPI != nil {
		snap.PPIID = &c.PPI.ID
	}
	return snap
}

type BuildContextRequest struct {
	OrgID     uuid.UUID
	UserID    int64
	Dimension string
	StudentID *int64
	Topic     string
}

func (r BuildContextRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.UserID <= 0 {
		return errUserIDRequired
	}
	return nil
}

type BuildPromptContext interface {
	Execute(ctx context.Context, req BuildContextRequest) (*PromptContext, error)
}

type buildPromptContextImpl struct {
	students    providers.StudentProvider
	teachers    providers.TeacherProfileProvider
	situations  providers.SituationCatalogProvider
	diagnoses   providers.DiagnosisProvider
	ppi         providers.PPIProvider
	adaptations providers.AdaptationProvider
	classrooms  providers.ClassroomProvider
	devices     providers.DeviceProvider
	summaries   providers.ConversationSummaryProvider
}

func NewBuildPromptContext(
	students providers.StudentProvider,
	teachers providers.TeacherProfileProvider,
	situations providers.SituationCatalogProvider,
	diagnoses providers.DiagnosisProvider,
	ppi providers.PPIProvider,
	adaptations providers.AdaptationProvider,
	classrooms providers.ClassroomProvider,
	devices providers.DeviceProvider,
	summaries providers.ConversationSummaryProvider,
) BuildPromptContext {
	return &buildPromptContextImpl{
		students:    students,
		teachers:    teachers,
		situations:  situations,
		diagnoses:   diagnoses,
		ppi:         ppi,
		adaptations: adaptations,
		classrooms:  classrooms,
		devices:     devices,
		summaries:   summaries,
	}
}

func (uc *buildPromptContextImpl) Execute(ctx context.Context, req BuildContextRequest) (*PromptContext, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	pc := &PromptContext{Dimension: req.Dimension}

	// ---- STATIC (eager): device catalog + situation vocabulary. ----
	devices, err := uc.devices.ListDevices(ctx, req.OrgID, nil)
	if err != nil {
		return nil, err
	}
	pc.DeviceCatalog = devices

	situations, err := uc.situations.List(ctx, req.OrgID)
	if err != nil {
		return nil, err
	}
	pc.Situations = situations

	// ---- Teacher (who we are talking to): optional, suggested if missing. ----
	teacher, err := uc.teachers.GetByUserID(ctx, req.OrgID, req.UserID)
	switch {
	case err == nil:
		pc.Teacher = teacher
	case errors.Is(err, providers.ErrNotFound):
		pc.MissingData = append(pc.MissingData, missingTeacherProfile)
	default:
		return nil, err
	}

	// ---- DYNAMIC lazy: only the requested dimension. ----
	if req.Dimension == DimensionStudent && req.StudentID != nil && *req.StudentID > 0 {
		if err := uc.loadStudentDimension(ctx, req, pc); err != nil {
			return nil, err
		}
	}

	return pc, nil
}

// loadStudentDimension fetches the focus-student context: profile + diagnoses +
// PPI + classroom + past adaptations + prior summaries. All optional — missing
// pieces are recorded in MissingData rather than causing an error.
func (uc *buildPromptContextImpl) loadStudentDimension(ctx context.Context, req BuildContextRequest, pc *PromptContext) error {
	student, err := uc.students.GetStudent(ctx, req.OrgID, *req.StudentID)
	if err != nil {
		return err
	}
	pc.TargetStudent = student

	if student.Profile != nil {
		diags, err := uc.diagnoses.ListByStudentProfile(ctx, student.Profile.ID)
		if err != nil {
			return err
		}
		pc.Diagnoses = diags
		if len(diags) == 0 {
			pc.MissingData = append(pc.MissingData, missingDiagnoses)
		}
	} else {
		pc.MissingData = append(pc.MissingData, missingStudentProfile)
	}

	ppi, err := uc.ppi.GetByStudentID(ctx, req.OrgID, *req.StudentID)
	switch {
	case err == nil:
		pc.PPI = ppi
	case errors.Is(err, providers.ErrNotFound):
		pc.MissingData = append(pc.MissingData, missingPPI)
	default:
		return err
	}

	adaptations, err := uc.adaptations.List(ctx, req.OrgID, providers.AdaptationFilter{StudentID: req.StudentID})
	if err != nil {
		return err
	}
	if len(adaptations) > maxPastAdaptations {
		adaptations = adaptations[:maxPastAdaptations]
	}
	pc.PastAdaptations = adaptations

	if student.ClassroomID > 0 {
		classroom, err := uc.classrooms.Get(ctx, req.OrgID, student.ClassroomID)
		switch {
		case err == nil:
			pc.Classroom = classroom
			peers, err := uc.students.ListByClassroom(ctx, req.OrgID, student.ClassroomID)
			if err != nil {
				return err
			}
			pc.ClassroomStudents = peers
		case errors.Is(err, providers.ErrNotFound):
			// classroom deleted or inconsistent — continue without it
		default:
			return err
		}
	}

	summaries, err := uc.summaries.RecentByStudent(ctx, req.OrgID, *req.StudentID, maxPriorSummaries)
	if err != nil {
		return err
	}
	pc.PriorSummaries = summaries

	return nil
}
