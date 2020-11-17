package main

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go" // https://godoc.org/firebase.google.com/go
)

// BT - BubbleTimer application
type BT struct {
	config          Config
	currentDay      Day
	ui              UI
	firebaseApp     *firebase.App
	firebaseContext context.Context
	firestoreClient *firestore.Client
}

var bt BT

func main() {
	bt.config = getConfig()
	bt.firebaseApp, bt.firebaseContext, bt.firestoreClient = firebaseConnect()
	bt.currentDay = loadData(time.Now())
	bt.ui = initUI()
	startUI() // Blocking
	shutdown()
}

func shutdown() {
	bt.firestoreClient.Close()
	fmt.Println("Come back soon!")
}
