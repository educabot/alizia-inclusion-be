package inclusion

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/educabot/alizia-inclusion-be/src/core/entities"
	"github.com/educabot/alizia-inclusion-be/src/core/providers"
)

// Dimensions a teacher can use to open a session. There is no mode selection;
// Alizia adapts its behaviour based on the dimension and the current context.
const (
	DimensionStudent = "alumno"
	DimensionToolkit = "valija"
	DimensionTopic   = "tema"
)

// maxPriorSummaries caps summary retrieval at session open to prevent context explosion.
const maxPriorSummaries = 10

const (
	welcomeText = "¡Hola! Soy Alizia, tu asistente de inclusión. ¿De qué querés hablar: " +
		"de un alumno, de la valija o de un tema?"
	clarifyText = "No me quedó claro. ¿Querés hablar de un alumno, de la valija o de un tema?"
	askStudent  = "Dale, ¿de qué alumno querés hablar?"
	askTopic    = "Perfecto, ¿sobre qué tema querés que busquemos?"
)

type OpenSessionRequest struct {
	OrgID     uuid.UUID
	UserID    int64
	Dimension string // alumno / valija / tema; empty = greeting only
	StudentID *int64 // required when Dimension = alumno
	DeviceID  *int64 // optional when Dimension = valija
	Topic     string // required when Dimension = tema
}

func (r OpenSessionRequest) Validate() error {
	if r.OrgID == uuid.Nil {
		return errOrgIDRequired
	}
	if r.UserID <= 0 {
		return errUserIDRequired
	}
	return nil
}

type OpenSessionResponse struct {
	Greeting       string                         `json:"greeting"`
	NeedsDimension bool                           `json:"needs_dimension"`
	Dimension      string                         `json:"dimension,omitempty"`
	Student        *entities.Student              `json:"student,omitempty"`
	PriorSummaries []entities.ConversationSummary `json:"prior_summaries,omitempty"`
}

type OpenSession interface {
	Execute(ctx context.Context, req OpenSessionRequest) (*OpenSessionResponse, error)
}

type openSessionImpl struct {
	students  providers.StudentProvider
	summaries providers.ConversationSummaryProvider
}

func NewOpenSession(students providers.StudentProvider, summaries providers.ConversationSummaryProvider) OpenSession {
	return &openSessionImpl{students: students, summaries: summaries}
}

func (uc *openSessionImpl) Execute(ctx context.Context, req OpenSessionRequest) (*OpenSessionResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	dimension := strings.ToLower(strings.TrimSpace(req.Dimension))

	switch dimension {
	case "":
		return &OpenSessionResponse{Greeting: welcomeText, NeedsDimension: true}, nil

	case DimensionStudent:
		return uc.openStudent(ctx, req)

	case DimensionToolkit:
		return uc.openToolkit(ctx, req)

	case DimensionTopic:
		return uc.openTopic(ctx, req)

	default:
		// Ambiguous dimension: ask again rather than assume.
		return &OpenSessionResponse{Greeting: clarifyText, NeedsDimension: true}, nil
	}
}

func (uc *openSessionImpl) openStudent(ctx context.Context, req OpenSessionRequest) (*OpenSessionResponse, error) {
	if req.StudentID == nil || *req.StudentID <= 0 {
		return &OpenSessionResponse{Greeting: askStudent, NeedsDimension: true}, nil
	}

	student, err := uc.students.GetStudent(ctx, req.OrgID, *req.StudentID)
	if err != nil {
		return nil, err
	}

	summaries, err := uc.summaries.RecentByStudent(ctx, req.OrgID, *req.StudentID, maxPriorSummaries)
	if err != nil {
		return nil, err
	}

	return &OpenSessionResponse{
		Greeting:       "Listo, hablemos de " + student.Name + ".",
		Dimension:      DimensionStudent,
		Student:        student,
		PriorSummaries: summaries,
	}, nil
}

func (uc *openSessionImpl) openToolkit(ctx context.Context, req OpenSessionRequest) (*OpenSessionResponse, error) {
	// The toolkit is injected as a catalogue in context, not via retrieval. Prior summaries
	// are fetched only when a specific device is targeted.
	var summaries []entities.ConversationSummary
	if req.DeviceID != nil && *req.DeviceID > 0 {
		var err error
		summaries, err = uc.summaries.RecentByDevice(ctx, req.OrgID, *req.DeviceID, maxPriorSummaries)
		if err != nil {
			return nil, err
		}
	}
	return &OpenSessionResponse{
		Greeting:       "Genial, miremos la valija. ¿Qué necesitás resolver?",
		Dimension:      DimensionToolkit,
		PriorSummaries: summaries,
	}, nil
}

func (uc *openSessionImpl) openTopic(ctx context.Context, req OpenSessionRequest) (*OpenSessionResponse, error) {
	topic := strings.TrimSpace(req.Topic)
	if topic == "" {
		return &OpenSessionResponse{Greeting: askTopic, NeedsDimension: true}, nil
	}

	summaries, err := uc.summaries.RecentByTopic(ctx, req.OrgID, topic, maxPriorSummaries)
	if err != nil {
		return nil, err
	}

	return &OpenSessionResponse{
		Greeting:       "Buenísimo, trabajemos sobre " + topic + ".",
		Dimension:      DimensionTopic,
		PriorSummaries: summaries,
	}, nil
}
