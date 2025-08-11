package service

import (
	"fmt"
	"time"
)

// type EnrichedVideoInfo  struct{
//   utility.APIVideoInfo
// }

// func (v *EnrichedVideoInfo) TimeToStart() string {
// 	if v.StartScheduled == "" {
// 		return ""
// 	}

// 	startTime, err := time.Parse(time.RFC3339, v.StartScheduled)
// 	if err != nil {
// 		// Try fallback without timezone
// 		startTime, err = time.ParseInLocation("2006-01-02T15:04:05", v.StartScheduled, time.UTC)
// 		if err != nil {
// 			// Log or skip invalid field
// 			fmt.Printf("Invalid start_scheduled value for video %s: %s\n", v.Title, v.StartScheduled)
// 			return ""
// 		}
// 	}

// 	diff := startTime.Sub(TimeNow())
// 	return fmt.Sprintf("in %s", FormatDuration(diff))
// }

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60

	if h > 0 && m > 0 {
		return fmt.Sprintf("%dh%dm", h, m)
	} else if h > 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dm", m)
}

// Simulate timeNow function
var TimeNow = func() time.Time {
	return time.Now().UTC()
}
