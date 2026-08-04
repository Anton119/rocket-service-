package auth

const (
	// SessionMetadataKey — ключ session_uuid в gRPC metadata и Kafka headers.
	SessionMetadataKey = "session-uuid"
	// SessionHeader — HTTP-заголовок с UUID сессии (альтернатива Bearer).
	SessionHeader = "X-Session-UUID"
)
