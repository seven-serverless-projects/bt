package main

import (
	"fmt"
	//"github.com/gdamore/tcell" // https://github.com/gdamore/tcell
	"github.com/rivo/tview" // https://github.com/rivo/tview
)

func main() {

	conf := getConfig()
	today := retrieveData(conf)
	InitUI(conf, today)
	RunUI(conf, today) // blocking
	shutdown()
}

func InitUI(conf Config, day Day) {
	app := tview.NewApplication()

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	list := tview.NewList().
		AddItem("7:15-7:00", "", '1', nil).
		AddItem("7:30-6:45", "", '2', nil).
		AddItem("7:45-8:00", "Writing", '3', nil).
		AddItem("8:00-8:15", "", '4', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	grid := tview.NewGrid().
		SetRows(3, 0, 3).
		AddItem(newPrimitive("Header"), 0, 0, 1, 3, 0, 0, false).
		AddItem(newPrimitive("Footer"), 2, 0, 1, 3, 0, 0, false).
		AddItem(list, 1, 0, 1, 3, 0, 0, false)

	if err := app.SetRoot(grid, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

func RunUI(conf Config, day Day) {

}

func shutdown() {
	fmt.Println("Come back soon!")
}
