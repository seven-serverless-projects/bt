package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go" // https://godoc.org/firebase.google.com/go
)

const dateFormat = "2006-01-02"

// Day - a single day of 96 (15m) time slices
type Day struct {
	date       string        // ISO 8601
	timeSlices [96]TimeSlice // 0 to 95 for each 15m of a day, 0 = 0:00-0:15, 4 = 1:00-1:15, 95 = 23:45-24:00
}

// TimeSlice - one unit of time, either uncategorized, or associated with at activity
type TimeSlice struct {
	slice      int
	activityID string
}

func firebaseConnect() (*firebase.App, context.Context, *firestore.Client) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: bt.config.ProjectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		fmt.Println("\nUnable to initialize Firebase.")
		panic(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		fmt.Println("\nUnable to connect to Firestore.")
		panic(err)
	}
	return app, ctx, client
}

// TODO for a different day than today
// TODO retrieve from Firestore rather than blank
//
func retrieveData() Day {
	day := Day{}
	day.date = time.Now().Format(dateFormat)
	for i, slice := range day.timeSlices {
		slice.slice = i
		day.timeSlices[i] = slice
	}
	return day
}

//
func currentTimeSlicesFor(day Day) []TimeSlice {
	now := time.Now()

	// time slices between now and the top of the current hour
	nowMinutes := now.Minute()
	thisHourSlices := (nowMinutes / 15) + 1

	// time slices before the top of the hour
	priorHourSlices := timeSlicesDisplayed - thisHourSlices

	// time slice index for this hour
	nowHour := now.Hour()
	startingTimeSlice := (nowHour * 4) - priorHourSlices

	// Adjust for being near the the start and end of the day
	if startingTimeSlice < 0 {
		startingTimeSlice = 0
	} else if startingTimeSlice > (96 - timeSlicesDisplayed) {
		startingTimeSlice = 96 - timeSlicesDisplayed
	}

	// select and return the time slices from day
	return day.timeSlices[startingTimeSlice:(startingTimeSlice + timeSlicesDisplayed)]
}

//
func activityByID(id string) Activity {
	var matchingActivity Activity
	for _, activity := range bt.config.Activities {
		if activity.ID == id {
			matchingActivity = activity
			break
		}
	}
	return matchingActivity
}

//
func activeActivities() []Activity {
	activeActivities := []Activity{}
	for _, activity := range bt.config.Activities {
		if activity.Active {
			activeActivities = append(activeActivities, activity)
		}
	}
	return activeActivities
}

//
func persist() (bool, string) {
	success := true
	errorMessage := ""

	_, err := bt.firestoreClient.
		Collection("users").
		Doc(bt.config.UserID).
		Collection("days").
		Doc(bt.currentDay.date).
		Set(bt.firebaseContext, sparseTimeSliceActivityMap(bt.currentDay.timeSlices[:]))

	if err != nil {
		success = false
		errorMessage = err.Error()
		fmt.Printf("Firestore write - Error: %s", errorMessage)
	}
	return success, errorMessage
}

//
func sparseTimeSliceActivityMap(timeSlices []TimeSlice) map[string]string {
	timeSliceMap := make(map[string]string)
	for i := range timeSlices {
		if timeSlices[i].activityID != "" {
			timeSliceMap[fmt.Sprint(i)] = timeSlices[i].activityID
		}
	}
	return timeSliceMap
}
