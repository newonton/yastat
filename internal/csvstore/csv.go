package csvstore

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/newonton/yastat/internal/timezone"
)

func LastTime(path string) (time.Time, bool) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, false
	}
	defer f.Close()

	r := csv.NewReader(f)
	var rec []string
	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		rec = row
	}
	if len(rec) == 0 {
		return time.Time{}, false
	}

	t, err := time.ParseInLocation(time.DateTime, rec[0], timezone.MoscowZone)
	if err != nil {
		return time.Time{}, false
	}

	return t.Truncate(time.Minute), true
}

func Append(path string, t time.Time, reward float64, shows int) error {
	f, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	return w.Write([]string{
		t.Format(time.DateTime),
		formatFloat(reward),
		formatInt(shows),
	})
}

func formatFloat(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }
func formatInt(v int) string       { return strconv.Itoa(v) }
