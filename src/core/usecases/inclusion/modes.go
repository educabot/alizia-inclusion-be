package inclusion

// Conversation modes tag each AI turn (persisted on the turn and on the usage
// trace) so history and analytics can tell apart how the teacher was assisted.
const (
	modeAssist    = "assist"
	modeRecommend = "recommend"
	modeClose     = "close"
)
