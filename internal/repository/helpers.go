package repository
import (
	"time"
)

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nilIfTimeZero(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
}

func ptrToTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
