package main

import (
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go" // https://godoc.org/firebase.google.com/go
)

// BT - BubbleTimer application
type BT struct {
	config          Config
	currentDay      Day
	ui              UI
	firebaseApp     *firebase.App
	firestoreClient *firestore.Client
}

var bt BT

func main() {
	bt.config = getConfig()
	bt.firebaseApp, bt.firestoreClient = firebaseConnect()
	bt.currentDay = retrieveData()
	bt.ui = initUI()
	startUI() // Blocking
	shutdown()
}

func shutdown() {
	fmt.Println("Come back soon!")
}
