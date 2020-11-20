package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2" // https://github.com/gdamore/tcell
	"github.com/rivo/tview"       // https://github.com/rivo/tview
)

const (
	timeSlicesDisplayed = 12
	bgColor             = tcell.ColorDarkBlue
	formatUS            = "Monday, January 2, 2006"
)

// UI - the BubbleTimer terminal user interface
type UI struct {
	app               *tview.Application
	grid              *tview.Grid
	header            *tview.TextView
	timeSliceList     *tview.TextView
	activityList      *tview.TextView
	commandInput      *tview.InputField
	currentTimeSlices []TimeSlice
}

var ui UI

func initUI() UI {
	ui.app = tview.NewApplication()

	initRegExp()
	initHeader()
	initTimeSlices()
	initActivities()
	initFooter()
	initGrid()

	return ui
}

func startUI() {
	if err := ui.app.Run(); err != nil {
		panic(err)
	}
}

func initHeader() {
	thisDay, err := time.Parse(dateFormat, bt.currentDay.date)
	if err != nil {
		fmt.Println("Unable to parse the date: " + bt.currentDay.date)
		panic(err)
	}
	ui.header = tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(thisDay.Format(formatUS))
	ui.header.SetBorderPadding(1, 1, 0, 0)
	ui.header.SetTextColor(tcell.ColorLimeGreen)
	ui.header.SetBackgroundColor(bgColor)
}

func initTimeSlices() {
	ui.timeSliceList = tview.NewTextView()
	ui.timeSliceList.SetBorderPadding(0, 0, 1, 1).
		SetBackgroundColor(bgColor)
	ui.currentTimeSlices = timeSlicesForTime(bt.currentDay, time.Now())
	ui.timeSliceList.SetText(timeSliceText())
}

func initActivities() {
	ui.activityList = tview.NewTextView()
	ui.activityList.SetBorderPadding(0, 0, 1, 1).
		SetBackgroundColor(bgColor)
	ui.activityList.SetText(activityText())
}

func initFooter() {
	ui.commandInput = tview.NewInputField().
		SetLabel("Command: ").
		SetFieldWidth(25).
		SetFieldBackgroundColor(bgColor).
		SetFieldTextColor(tcell.ColorYellow).
		SetLabelColor(tcell.ColorGreen).
		SetDoneFunc(inputComplete).
		SetText("")
	ui.commandInput.SetBorderPadding(1, 1, 1, 1)
	ui.commandInput.SetBackgroundColor(bgColor)
}

func initGrid() {
	ui.grid = tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(0, 0).
		AddItem(ui.header, 0, 0, 1, 2, 0, 0, false).
		AddItem(ui.timeSliceList, 1, 0, 1, 1, 0, 0, false).
		AddItem(ui.activityList, 1, 1, 1, 1, 0, 0, false).
		AddItem(ui.commandInput, 2, 0, 1, 2, 0, 0, true)
	ui.app.SetRoot(ui.grid, true)
	ui.app.SetFocus(ui.commandInput)
}

func resetInput() {
	if ui.commandInput != nil {
		ui.commandInput.SetText("")
	}
}

// Return a string suitable for use in the UI with a return delimitted entry for each timeslice we are displaying
func timeSliceText() string {
	timeSliceText := ""
	for i := range ui.currentTimeSlices {
		timeSlice := ui.currentTimeSlices[i]
		timeSliceText += "t" + fmt.Sprint(i+1) + " — " + timeDisplayFor(timeSlice)
		if timeSlice.activityID != "" {
			activity := activityByID(timeSlice.activityID)
			if activity.Name != "" {
				timeSliceText += " — " + activity.Name
			}
		}
		timeSliceText += "\n\n"
	}
	return timeSliceText
}

// Return a string suitable for use in the UI with a return delimitted entry for each activity we are displaying
// TODO Right now this just uses active activities, but the specified day may have older, inactive activities as well
func activityText() string {
	activeActivityCount := 1
	activityText := ""
	for _, activity := range activeActivities() {
		activityText += "a" + fmt.Sprint(activeActivityCount) + " — " + activity.Name
		timeInActivity := timeInActivityText(activity.ID)
		if timeInActivity != "" {
			activityText += " — " + timeInActivity
		}
		activityText += "\n\n"
		activeActivityCount++
	}
	return activityText
}

// Takes a time slice and returns a human readable string representing the starting and ending time of the time slice.
// Currently in 24h time only.
func timeDisplayFor(timeSlice TimeSlice) string {
	startHour := timeSlice.slice / 4
	startMinute := (timeSlice.slice % 4)
	endHour := startHour
	if startMinute == 3 {
		endHour = startHour + 1
	}
	startMinuteString := fmt.Sprint(startMinute * 15)
	if startMinute == 0 {
		startMinuteString = "00"
	}
	endMinute := ((timeSlice.slice % 4) + 1)
	if endMinute == 4 {
		endMinute = 0
	}
	endMinuteString := fmt.Sprint(endMinute * 15)
	if endMinute == 0 {
		endMinuteString = "00"
	}
	return fmt.Sprintf("%d:%s - %d:%s %s", startHour, startMinuteString, endHour, endMinuteString, "")
}

// Sum any timeslices spent doing the specified activity during the displayed day
// into human readable text e.g. 2h 15m
// Return a blank string if there's no timeslices for the activity.
func timeInActivityText(activityID string) string {
	timeInActivity := ""
	timeSliceCount := 0
	for _, timeslice := range bt.currentDay.timeSlices {
		if timeslice.activityID == activityID {
			timeSliceCount++
		}
	}
	if timeSliceCount > 0 {
		hours := timeSliceCount / 4
		minuteSlices := timeSliceCount % 4
		if hours > 0 {
			timeInActivity = fmt.Sprintf("%dh ", hours)
		}
		if minuteSlices > 0 {
			var minutes int
			switch minuteSlices {
			case 1:
				minutes = 15
			case 2:
				minutes = 30
			case 3:
				minutes = 45
			}
			timeInActivity += fmt.Sprintf("%dm", minutes)
		}
	}
	return timeInActivity
}

// The user finished their input, if they finished it with enter, attempt to parse it, otherwise reset the input
func inputComplete(key tcell.Key) {
	if key == tcell.KeyEnter {
		parseInput()
	} else {
		resetInput()
	}
}

// Assign the specified activity to the specified time slices and persist the update
func assignTime(timeSliceIndexes []int, activityIndex int) {

	// Get the activity
	activity := activeActivities()[activityIndex-1]

	// update the day's timeslices with the activity
	for _, timeSliceIndex := range timeSliceIndexes {
		timeSlice := ui.currentTimeSlices[timeSliceIndex-1]
		timeSlice.activityID = activity.ID // set the activity
		// Replace the time slice in the UI's data
		ui.currentTimeSlices[timeSliceIndex-1] = timeSlice
		// Replace the time slice in the current day's data
		timeSlices := bt.currentDay.timeSlices
		timeSlices[timeSlice.slice] = timeSlice
		bt.currentDay.timeSlices = timeSlices
	}

	// refresh the timeslices and activity display in the ui
	ui.timeSliceList.SetText(timeSliceText())
	ui.activityList.SetText(activityText())

	// persist the updated timeslices
	// TODO status message
	persist()
	// TODO status message
}

func unassignTime(timeSliceIndexes []int) {
	fmt.Printf("Unassign: %v", timeSliceIndexes)
}

// Increment startingTimeSlice by page size (adjusting for end of day) and rerender
func timeForward() {
	// start with the last of the current time slices
	newStartingTimeSlice := ui.currentTimeSlices[len(ui.currentTimeSlices)-1]
	ui.currentTimeSlices = timeSlicesForIndex(bt.currentDay, newStartingTimeSlice.slice)
	ui.timeSliceList.SetText(timeSliceText())
}

// Decrement startingTimeSlice by page size (adjusting for start of day) and rerender
func timeBackward() {
	// end with the first of the current time slices
	newEndingTimeSlice := ui.currentTimeSlices[0]
	ui.currentTimeSlices = timeSlicesForIndex(bt.currentDay, newEndingTimeSlice.slice-timeSlicesDisplayed+1)
	ui.timeSliceList.SetText(timeSliceText())
}

// Parse the current day, increment it by one, and reset the UI
func dayForward() {
	current, err := time.Parse(dateFormat, bt.currentDay.date)
	if err != nil {
		fmt.Println("Unable to parse current day: " + bt.currentDay.date)
		panic(err)
	}
	forward := current.Add(-time.Hour * 24)
	resetForDay(forward)
}

// Parse the current day, decrement it by one, and reset the UI
func dayBackward() {
	current, err := time.Parse(dateFormat, bt.currentDay.date)
	if err != nil {
		fmt.Println("Unable to parse current day: " + bt.currentDay.date)
		panic(err)
	}
	back := current.Add(-time.Hour * 24)
	resetForDay(back)
}

func dayTodayTimeNow() {
	now := time.Now()
	// Set the timeslices to end at the current time
	ui.currentTimeSlices = timeSlicesForTime(bt.currentDay, now)
	// Set the day to today
	resetForDay(now)
}

func dayYesterday() {
	now := time.Now()
	yesterday := now.Add(-time.Hour * 24)
	resetForDay(yesterday)
}

func resetForDay(day time.Time) {
	bt.currentDay = loadData(day)
	// Reset the UI, using the same starting time slice as is currently shown
	ui.currentTimeSlices = timeSlicesForIndex(bt.currentDay, ui.currentTimeSlices[0].slice)
	ui.header.SetText(day.Format(formatUS))
	ui.timeSliceList.SetText(timeSliceText())
	ui.activityList.SetText(activityText())
	// TODO account for any activities that are on the day but not active
}
