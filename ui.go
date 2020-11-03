package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2" // https://github.com/gdamore/tcell
	"github.com/rivo/tview"       // https://github.com/rivo/tview
)

const (
	bgColor  = tcell.ColorDarkBlue
	formatUS = "Monday, November 2, 2020"
)

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

func initUI() UI {
	ui.app = tview.NewApplication()

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
		SetFieldWidth(10).
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
			activityText += "c" + fmt.Sprint(activeActivityCount) + " — " + activity.Name + "\n"
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

func parseInput() {
	input := ui.commandInput.GetText()
	if input == "q" || input == "quit" || input == "Q" || input == "Quit" {
		ui.app.Stop()
	} else {
		resetInput()
	}
}

func resetInput() {
	ui.commandInput.SetText("")
}
