package models

import (
	"fmt"
	"time"
)

func timeAgo(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	since := time.Since(t)
	switch {
	case since < time.Minute:
		return "just now"
	case since < time.Hour:
		return formatTimeText(int(since.Minutes()), "minute")
	case since < 48*time.Hour:
		return formatTimeText(int(since.Hours()), "hour")
	default:
		return formatTimeText(int(since.Hours()/24), "day")
	}
}

func formatTimeText(n int, unit string) string {
	if n == 1 {
		return "1 " + unit + " ago"
	}
	return fmt.Sprintf("%d %ss ago", n, unit)
}
