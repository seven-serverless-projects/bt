package main

import (
	"fmt"
)

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
	shutdown()
}

func shutdown() {
	fmt.Println("Come back soon!")
}
