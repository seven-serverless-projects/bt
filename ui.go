package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gdamore/tcell/v2" // https://github.com/gdamore/tcell
	"github.com/rivo/tview"       // https://github.com/rivo/tview
)

const (
	bgColor  = tcell.ColorDarkBlue
	formatUS = "Monday, November 2, 2020"
)

/*
Regular expression that can parse valid time entries from the user.

t# references the time slices displayed in the UI by their numbered index.

a# references activities displayed in the UI by their numbered index.

One or more time slices, or a range of time slices, are associated with one activity.

Valid Examples:
t1 a1
t1a1
t3, t6 a2
t3,t6 a2
t3 t6 a2
t3t6 a2
t3t6a2
t7-t10 a5
t7-t10a5

Detailed breakdown of the regex string:

^ - the start

(?:t(?P<sliceIndex>[0-9]+),?\\s?)* - a single t# or a comma, space, or non-delimitted sequence of them

| - or this other way of specifiying time

t(?P<range1>[0-9]+)-t(?P<range2>[0-9]+) - a range t#=t#

\s optional white space between time and activity

a(?P<activity>[0-9]+) - activity in the form of a#

$ - the end
*/
const timeEntryRegExString = "^((?:t(?P<sliceIndex>[0-9]+),?\\s?)*|t(?P<range1>[0-9]+)-t(?P<range2>[0-9]+))\\s*a(?P<activity>[0-9]+)$"

// Subset of the full regex, with just a single t# or a comma, space, or non-delimitted sequence of them
const timeSlicesRegExString = "^(t([0-9])+,?\\s?)+$"

// UI - the BubbleTimer terminal user interface
type UI struct {
	app           *tview.Application
	grid          *tview.Grid
	header        *tview.TextView
	timeSliceList *tview.TextView
	activityList  *tview.TextView
	commandInput  *tview.InputField
}

var ui UI
var timeEntryRegExp, timeSlicesRegExp *regexp.Regexp

func initUI() UI {
	ui.app = tview.NewApplication()

	initRegExp()
	initHeader()
	initTimeSlices()
	initActivities()
	initFooter()
	initGrid()

	if err := ui.app.SetRoot(ui.grid, true).Run(); err != nil {
		fmt.Println("Unable to initialize the UI!")
		panic(err)
	}

	return ui
}

func initRegExp() {
	timeEntryRegExp = regexp.MustCompile(timeEntryRegExString)
	timeSlicesRegExp = regexp.MustCompile(timeSlicesRegExString)
}

func initHeader() {
	thisDay, err := time.Parse(time.RFC3339, bt.currentDay.date)
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
	setTimeSliceTextFor(bt.currentDay)
}

func initActivities() {
	ui.activityList = tview.NewTextView()
	ui.activityList.SetBorderPadding(0, 0, 1, 1).
		SetBackgroundColor(bgColor)
	setActivitiesFor(bt.currentDay)
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
}

func resetInput() {
	if ui.commandInput != nil {
		ui.commandInput.SetText("")
	}
}

func setTimeSliceTextFor(thisDay Day) {
	currentTimeSlices := currentTimeSlicesFor(thisDay)
	timeSliceText := ""
	for i := range currentTimeSlices {
		k := len(currentTimeSlices) - i - 1 // backwards iteration
		timeSlice := currentTimeSlices[k]
		timeSliceText += "t" + fmt.Sprint(i+1) + " — " + timeDisplayFor(timeSlice)
		if timeSlice.activityID != "" {
			activity := activityByID(timeSlice.activityID)
			if activity.Name != "" {
				timeSliceText += " — " + activity.Name
			}
		}
		timeSliceText += "\n"
	}
	ui.timeSliceList.SetText(timeSliceText)
}

func setActivitiesFor(thisDay Day) {
	activeActivityCount := 1
	activityText := ""
	for _, activity := range activeActivities() {
		activityText += "a" + fmt.Sprint(activeActivityCount) + " — " + activity.Name + "\n"
		activeActivityCount++
	}
	ui.activityList.SetText(activityText)
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

func inputComplete(key tcell.Key) {
	if key == tcell.KeyEnter {
		parseInput()
	} else {
		resetInput()
	}
}

func assignTime(timeSlices []int, activity int) {

	// Get the index offset of the current block of time slices

	// update the day's timeslices with the activity

	// refresh the timeslices display in the ui
	initTimeSlices()
}
