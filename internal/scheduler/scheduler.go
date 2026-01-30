package scheduler

import (
	"time"

	"github.com/newonton/yastat/internal/timezone"
)

func NextRun(last time.Time, period time.Duration) time.Time {
	return last.In(timezone.MoscowZone).Truncate(time.Minute).Add(period)
}
