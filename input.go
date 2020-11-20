package main

import (
	"regexp"
	"strconv"
	"strings"
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
t7-10 a5
t7-t10a5
t7-10a5

Detailed breakdown of the regex string:

^ - the start

(?:t(?P<sliceIndex>[0-9]+),?\\s?)* - a single t# or a comma, space, or non-delimitted sequence of them

| - or this other way of specifiying time

t(?P<range1>[0-9]+)-t(?P<range2>[0-9]+) - a range t#=t#

\s optional white space between time and activity

a(?P<activity>[0-9]+) - activity in the form of a#

$ - the end
*/
const timeEntryRegExString = "^((?:t(?P<sliceIndex>[0-9]+),?\\s?)*|t(?P<range1>[0-9]+)-t?(?P<range2>[0-9]+))\\s*a(?P<activity>[0-9]+)$"

// Subset of the full parsing regex, with just a single t# or a comma, space, or non-delimitted sequence of them
const timeSlicesRegExString = "^(t([0-9])+,?\\s?)+$"

/*
Regular expression that can parse valid time unassignments from the user.

t# references the time slices displayed in the UI by their numbered index.

One or more time slices, or a range of time slices, are associated with one activity.

Valid Examples:
u t1
u t3, t6
u t3,t6
u t3 t6
u t3t6
ut3t6
u t7-t10
ut7-10

Detailed breakdown of the regex string:

^ - the start

u\\s* - the literal letter u and optional white space

(?:t(?P<sliceIndex>[0-9]+),?\\s?)* - a single t# or a comma, space, or non-delimitted sequence of them

| - or this other way of specifiying time

t(?P<range1>[0-9]+)-t(?P<range2>[0-9]+) - a range t#=t#

$ - the end
*/
const unassignRegExString = "^u\\s*((?:t(?P<sliceIndex>[0-9]+),?\\s?)*|t(?P<range1>[0-9]+)-t?(?P<range2>[0-9]+))$"

var timeEntryRegExp, timeSlicesRegExp, unassignRegExp *regexp.Regexp

func initRegExp() {
	timeEntryRegExp = regexp.MustCompile(timeEntryRegExString)
	timeSlicesRegExp = regexp.MustCompile(timeSlicesRegExString)
	unassignRegExp = regexp.MustCompile(unassignRegExString)
}

// Parse text input from the user, and do the requested action
func parseInput() {
	input := strings.ToLower(ui.commandInput.GetText())
	switch input {
	case "q", "quit":
		ui.app.Stop()
	case "+":
		timeForward()
	case "-":
		timeBackward()
	case "n", "next":
		dayForward()
	case "p", "prior":
		dayBackward()
	case "t", "today", "r", "refresh", "reset":
		dayTodayTimeNow()
	case "y", "yesterday":
		dayYesterday()
	default:
		if strings.HasPrefix(input, "t") {
			timeSlices, timeRange, activity, err := parseTimeEntry(input)
			if !err {
				if timeRange[0] > 0 {
					timeSlices = expandRange(timeRange[0], timeRange[1])
				}
				assignTime(timeSlices, activity)
			}
		} else if strings.HasPrefix(input, "u") {
			timeSlices, timeRange, err := parseUnassignment(input)
			if !err {
				if timeRange[0] > 0 {
					timeSlices = expandRange(timeRange[0], timeRange[1])
				}
				unassignTime(timeSlices)
			}
		}
	}
	resetInput()
}

// Parse the unassignment input from a user as integers.
// Time entry unassignment disassociates a time slice, multiple time slices,
// or a range of time slices from any assigned activity.
func parseUnassignment(entry string) ([]int, [2]int, bool) {
	timeSlices := []int{}
	var timeRange [2]int
	err := false
	matches := unassignRegExp.FindStringSubmatch(entry)
	if len(matches) < 5 ||
		!validRange(matches[3], matches[4]) {
		err = true // silly user!
	} else {
		timeSlices = parseTimeSlices(matches[1])
		timeRange[0], _ = strconv.Atoi(matches[3])
		timeRange[1], _ = strconv.Atoi(matches[4])
	}
	return timeSlices, timeRange, err
}

// Parse the time entry input from a user as integers.
// Time entry input associates a time slice, multiple time slices,
// or a range of time slices with a single activity.
func parseTimeEntry(entry string) ([]int, [2]int, int, bool) {
	timeSlices := []int{}
	var timeRange [2]int
	var activity int
	err := false
	matches := timeEntryRegExp.FindStringSubmatch(entry)
	if len(matches) < 6 ||
		!validRange(matches[3], matches[4]) ||
		!validActivity(matches[5]) {
		err = true // silly user!
	} else {
		timeSlices = parseTimeSlices(matches[1])
		timeRange[0], _ = strconv.Atoi(matches[3])
		timeRange[1], _ = strconv.Atoi(matches[4])
		activity, _ = strconv.Atoi(matches[5])
	}
	return timeSlices, timeRange, activity, err
}

// If the time portion of the user entry was provided as a time slice
// or as a sequence of time slices, than parse the integers of each
// provided time slice
func parseTimeSlices(entry string) []int {
	timeSlices := []int{}
	// sanity check to make sure we are dealing with a set of time slices
	if timeSlicesRegExp.MatchString(entry) {
		// Strip out white space and comma delimiters
		delimit := regexp.MustCompile("\\s*,*")
		compactEntry := delimit.ReplaceAllLiteralString(entry, "")
		// Slit on the t's
		t := regexp.MustCompile("t")
		timeSliceStrings := t.Split(compactEntry, -1)
		for i := 1; i < len(timeSliceStrings); i++ { // first item is always blank
			timeSlice, _ := strconv.Atoi(timeSliceStrings[i])
			timeSlices = append(timeSlices, timeSlice)
		}
	}
	return timeSlices
}

// Return true if the strings represent the start and end of a valid range of time slices
func validRange(startString string, endString string) bool {
	valid := true
	if startString != "" {
		start, _ := strconv.Atoi(startString)
		end, _ := strconv.Atoi(endString)
		if start <= 0 ||
			end <= 0 ||
			start >= end ||
			end > timeSlicesDisplayed {
			valid = false
		}
	}
	return valid
}

// Return true if the string represents the index of a valid activity from the UI
func validActivity(activity string) bool {
	valid := true
	index, _ := strconv.Atoi(activity) // safe due to the regexes
	if index < 1 || index > len(activeActivities()) {
		valid = false
	}
	return valid
}

// Given a range specified by 2 positive ints, the second bigger than the first, return an array of
// all the ints between them, inclusive of the start and end
func expandRange(start int, end int) []int {
	expanded := []int{}
	for i := start; i <= end; i++ {
		expanded = append(expanded, i)
	}
	return expanded
}
