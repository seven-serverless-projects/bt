package main

import (
	"fmt"
)

// BT - BubbleTimer application
type BT struct {
	config     Config
	currentDay Day
	ui         UI
}

var bt BT

func main() {

	bt.config = getConfig()
	bt.currentDay = retrieveData()
	bt.ui = initUI()
	startUI() // Blocking
	shutdown()
}

func shutdown() {
	fmt.Println("Come back soon!")
}
