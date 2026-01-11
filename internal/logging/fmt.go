package logging

// StringPtr returns a nil-safe string representation for logging.
// If the pointer is nil, it returns "<nil>". Otherwise, it returns the string value.
// This is useful for safely logging optional string values.
func StringPtr(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
