package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go" // https://godoc.org/firebase.google.com/go
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// Return a Firestore client that's connected to the app and ready to use
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

// Return an initialized day for the specified date.
// Load the document from Firestore for the specified day.
// Include any stored timeslice activities for the day in the data.
func loadData(forDay time.Time) Day {
	timeSliceMap := make(map[string]interface{})
	day := Day{}
	day.date = forDay.Format(dateFormat)
	doc, err := bt.firestoreClient.
		Collection("users").
		Doc(bt.config.UserID).
		Collection("days").
		Doc(day.date).
		Get(bt.firebaseContext)
	if (err != nil && status.Code(err) != codes.NotFound) || doc == nil {
		fmt.Printf("\nUnable to read data for: %s\n", day.date)
		panic(err)
	} else {
		timeSliceMap = doc.Data()
	}
	for i, slice := range day.timeSlices {
		slice.slice = i
		loadedData := timeSliceMap[fmt.Sprint(i)]
		if loadedData != nil { // the loaded time slice map is sparse
			activityID := loadedData.(map[string]interface{})["activity_id"]
			if activityID != nil {
				slice.activityID = activityID.(string) // Type conversion
			}
		}
		day.timeSlices[i] = slice
	}
	return day
}

// Return the configured number of time slices for the specified day,
// starting at the current time and working backwards in time (unless
// that takes us to midnight, in which case, use midnight as the earliest
// time slice).
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

// Given the ID of an activity, return the struct for the activity, returns
// nil if there's no match (not expected)
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

// Return an array of only the activities from the config that are active
func activeActivities() []Activity {
	activeActivities := []Activity{}
	for _, activity := range bt.config.Activities {
		if activity.Active {
			activeActivities = append(activeActivities, activity)
		}
	}
	return activeActivities
}

// Persist the timeslices for the current day being shown in the UI
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

// Given an array of all the timeslices for a day, create a map of just the timeslices
// with an assigned activity ID, using the timeslice index as the key
func sparseTimeSliceActivityMap(timeSlices []TimeSlice) map[string]map[string]string {
	timeSliceMap := make(map[string]map[string]string)
	for i := range timeSlices {
		if timeSlices[i].activityID != "" {
			timeSliceMap[fmt.Sprint(i)] = map[string]string{"activity_id": timeSlices[i].activityID}
		}
	}
	return timeSliceMap
}
