package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Parse text input from the user, which is request to quit, or a time entry.
func parseInput() {
	input := strings.ToLower(ui.commandInput.GetText())
	if input == "q" || input == "quit" {
		ui.app.Stop()
	} else if strings.HasPrefix(input, "t") {
		timeSlices, timeRange, activity, err := parseTimeEntry(input)
		if !err {
			if timeRange[0] > 0 {
				timeSlices = expandRange(timeRange[0], timeRange[1])
			}
			assignTime(timeSlices, activity)
		}
		resetInput()
	} else {
		resetInput()
	}
}

// Parse the time entry input from a user as integers.
// Time entry input associates a time slice, multiple time slices,
// or a range of time slices with a single activity
func parseTimeEntry(entry string) ([]int, [2]int, int, bool) {
	timeSlices := []int{}
	var timeRange [2]int
	var activity int
	err := false
	matches := timeEntryRegExp.FindStringSubmatch(entry)
	if len(matches) < 6 || !validRange(matches[3], matches[4]) {
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

// Given a range specified by 2 positive ints, the second bigger than the first, return an array of
// all the ints between them, inclusive of the start and end
func expandRange(start int, end int) []int {
	expanded := []int{}
	for i := start; i <= end; i++ {
		expanded = append(expanded, i)
	}
	return expanded
}
