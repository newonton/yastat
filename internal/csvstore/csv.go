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

type Record struct {
	Time        time.Time
	Shows       int
	Reward      float64
	DeltaShows  int
	DeltaReward float64
}

func LastRecord(path string) (Record, bool) {
	f, err := os.Open(path)
	if err != nil {
		return Record{}, false
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
		return Record{}, false
	}

	t, _ := time.ParseInLocation(time.DateTime, rec[0], timezone.MoscowZone)
	shows, _ := strconv.Atoi(rec[1])
	reward, _ := strconv.ParseFloat(rec[2], 64)

	var deltaShows int
	var deltaReward float64
	if len(rec) >= 5 {
		deltaShows, _ = strconv.Atoi(rec[3])
		deltaReward, _ = strconv.ParseFloat(rec[4], 64)
	}

	return Record{
		Time:        t.Truncate(time.Minute),
		Shows:       shows,
		Reward:      reward,
		DeltaShows:  deltaShows,
		DeltaReward: deltaReward,
	}, true
}

func Append(path string, t time.Time, shows int, reward float64, deltaShows int, deltaReward float64) error {
	f, _ := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	return w.Write([]string{
		t.In(timezone.MoscowZone).Format(time.DateTime),
		strconv.Itoa(shows),
		formatFloat(reward),
		strconv.Itoa(deltaShows),
		formatFloat(deltaReward),
	})
}

func formatFloat(v float64) string { return strconv.FormatFloat(v, 'f', 2, 64) }
func formatInt(v int) string       { return strconv.Itoa(v) }
