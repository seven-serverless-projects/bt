package main

import (
	"time"
)

type Day struct {
	Date       string        // ISO 8601
	TimeSlices [96]TimeSlice // 0 to 95 for each 15m of a day, 0 = 0:00-0:15, 4 = 1:00-1:15, 95 = 23:45-24:00
}

type TimeSlice struct {
	Slice          int
	TimeCategoryID string
}

// TODO for a different day than today
// TODO retrieve from Firestore rather than blank
func retrieveData(conf Config) Day {
	day := Day{}
	day.Date = time.Now().Format(time.RFC3339)
	for i, slice := range day.TimeSlices {
		slice.Slice = i
		day.TimeSlices[i] = slice
	}
	return day
}
