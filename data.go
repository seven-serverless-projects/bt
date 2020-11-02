package main

import (
	"time"
)

type Day struct {
	date       string        // ISO 8601
	timeSlices [96]TimeSlice // 0 to 95 for each 15m of a day, 0 = 0:00-0:15, 4 = 1:00-1:15, 95 = 23:45-24:00
}

type TimeSlice struct {
	slice          int
	timeCategoryID string
}

// TODO for a different day than today
// TODO retrieve from Firestore rather than blank
func retrieveData() Day {
	day := Day{}
	day.date = time.Now().Format(time.RFC3339)
	for i, slice := range day.timeSlices {
		slice.slice = i
		day.timeSlices[i] = slice
	}
	return day
}

func currentTimeSlicesFor(day Day) []TimeSlice {
	now := time.Now()

	// time slices between now and the top of the current hour
	nowMinutes := now.Minute()
	thisHourSlices := (nowMinutes / 15) + 1

	// time slices before the top of the hour
	priorHourSlices := 8 - thisHourSlices

	// time slice index for this hour
	nowHour := now.Hour()
	startingTimeSlice := (nowHour * 4) - priorHourSlices

	// Adjust for being near the the start and end of the day
	if startingTimeSlice < 0 {
		startingTimeSlice = 0
	} else if startingTimeSlice > 88 {
		startingTimeSlice = 88
	}

	// select and return the time slices from day
	return day.timeSlices[startingTimeSlice:(startingTimeSlice + 8)]
}
