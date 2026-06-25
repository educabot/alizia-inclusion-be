package inclusion

// Conversation-turn metadata keys. They are written when a turn is persisted
// (assist / recommend) and read back when a session is compacted (close), so a
// single source of truth prevents a silent mismatch between writer and reader.
const (
	metaKeyIdentifiedStudent = "identified_student"
	metaKeyRecommendedDevice = "recommended_device"
	metaKeyAdaptation        = "adaptation"
	metaKeySubject           = "subject"
)
