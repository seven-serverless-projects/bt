package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
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
var timeEntryRegExp *regexp.Regexp

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

func setTimeSliceTextFor(thisDay Day) {
	currentTimeSlices := currentTimeSlicesFor(thisDay)
	timeSliceText := ""
	for i := range currentTimeSlices {
		k := len(currentTimeSlices) - i - 1 // backwards iteration
		timeSlice := currentTimeSlices[k]
		timeSliceText += "t" + fmt.Sprint(i+1) + " — " + timeDisplayFor(timeSlice) + "\n"
	}
	ui.timeSliceList.SetText(timeSliceText)
}

func setActivitiesFor(thisDay Day) {
	activeActivityCount := 1
	activityText := ""
	for _, activity := range bt.config.Activities {
		if activity.Active {
			activityText += "a" + fmt.Sprint(activeActivityCount) + " — " + activity.Name + "\n"
			activeActivityCount++
		}
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

// Parse text input from the user, which is request to quit, or a time entry.
func parseInput() {
	input := strings.ToLower(ui.commandInput.GetText())
	if input == "q" || input == "quit" {
		ui.app.Stop()
	} else if strings.HasPrefix(input, "t") {
		parseTimeEntry(input)
	} else {
		resetInput()
	}
}

// The time entry format associates a single activity with a time slice, time slices, or a range of time slices
func parseTimeEntry(entry string) ([]int, [2]int, int, bool) {
	timeSlices := []int{}
	var timeRange [2]int
	var activity int
	err := false
	match := timeEntryRegExp.FindStringSubmatch(entry)
	if len(match) < 6 {
		err = true
		resetInput() // silly user!
	} else {
		resetInput()
		timeRange[0], _ = strconv.Atoi(match[3])
		timeRange[1], _ = strconv.Atoi(match[4])
		activity, _ = strconv.Atoi(match[5])
		//ui.commandInput.SetLabel("0: " + match[0] + " 1: " + match[1] + " 2: " + match[2] + " 3: " + match[3] + " 4: " + match[4] + " 5: " + match[5])
	}
	return timeSlices, timeRange, activity, err
}

func resetInput() {
	if ui.commandInput != nil {
		ui.commandInput.SetText("")
	}
}
