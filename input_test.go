package main

import (
	"fmt"
	"reflect"
	"testing"
)

const parseSuccess = false
const parseFailure = true

// TestParseTimeEntry - test user intput of a time to activity entry
func TestParseTimeEntry(t *testing.T) {

	type testCase struct {
		timeEntry  string
		timeSlices []int
		timeRange  [2]int
		activity   int
		err        bool
	}

	// "Stub" activeActivites func to return 5 (blank) activities
	active := Activity{"", "", "", parseFailure}
	bt.config.Activities = []Activity{active, active, active, active, active}

	testCases := []testCase{
		// failure cases - invalid time entries
		{"t1", []int{}, [2]int{}, 0, parseFailure},
		{"t1 a", []int{}, [2]int{}, 0, parseFailure},
		{"t1a", []int{}, [2]int{}, 0, parseFailure},
		{"t a1", []int{}, [2]int{}, 0, parseFailure},
		{"ta1", []int{}, [2]int{}, 0, parseFailure},
		{"ta", []int{}, [2]int{}, 0, parseFailure},
		{"ta a1", []int{}, [2]int{}, 0, parseFailure},
		{"ta1 a1", []int{}, [2]int{}, 0, parseFailure},
		{"t1 ab", []int{}, [2]int{}, 0, parseFailure},
		{"t1 ab1", []int{}, [2]int{}, 0, parseFailure},
		{"t1-t1 a1", []int{}, [2]int{}, 0, parseFailure},
		{"t0-t1 a1", []int{}, [2]int{}, 0, parseFailure},
		{"t2-t1 a1", []int{}, [2]int{}, 0, parseFailure},
		{"t1-t" + fmt.Sprint(timeSlicesDisplayed+1) + " a1", []int{}, [2]int{}, 0, parseFailure},
		// success cases - valid time entries
		{"t1 a1", []int{1}, [2]int{}, 1, parseSuccess},
		{"t1a1", []int{1}, [2]int{}, 1, parseSuccess},
		{"t3, t6 a2", []int{3, 6}, [2]int{}, 2, parseSuccess},
		{"t3,t6 a2", []int{3, 6}, [2]int{}, 2, parseSuccess},
		{"t3 t6 a2", []int{3, 6}, [2]int{}, 2, parseSuccess},
		{"t3t6 a2", []int{3, 6}, [2]int{}, 2, parseSuccess},
		{"t3t6a2", []int{3, 6}, [2]int{}, 2, parseSuccess},
		{"t7-t10 a5", []int{}, [2]int{7, 10}, 5, parseSuccess},
		{"t7-10 a5", []int{}, [2]int{7, 10}, 5, parseSuccess},
		{"t7-t10a5", []int{}, [2]int{7, 10}, 5, parseSuccess},
		{"t7-10a5", []int{}, [2]int{7, 10}, 5, parseSuccess}}
	t.Log("Test: parsing user time entries...")
	initRegExp()
	for i, testCase := range testCases {
		timeSlices, timeRange, activity, err := parseTimeEntry(testCase.timeEntry)
		if !err == testCase.err {
			t.Errorf("Test: parse entry FAIL -  parse outcome")
		} else if !reflect.DeepEqual(timeSlices, testCase.timeSlices) {
			t.Errorf("Test: parse entry FAIL - time slices in test case %d", i+1)
		} else if timeRange != testCase.timeRange {
			t.Errorf("Test: parse entry FAIL - time range in test case %d", i+1)
		} else if activity != testCase.activity {
			t.Errorf("Test: parse entry FAIL - activity in test case %d", i+1)
		} else {
			t.Log("Test: success for entry test case " + fmt.Sprint(i+1))
		}
	}
}

func TestParseUnassignment(t *testing.T) {

	type testCase struct {
		timeUnassignment string
		timeSlices       []int
		timeRange        [2]int
		err              bool
	}

	// "Stub" activeActivites func to return 5 (blank) activities
	active := Activity{"", "", "", parseFailure}
	bt.config.Activities = []Activity{active, active, active, active, active}

	testCases := []testCase{
		// failure cases
		{"u", []int{}, [2]int{}, parseFailure},
		{"u t", []int{}, [2]int{}, parseFailure},
		{"u ta", []int{}, [2]int{}, parseFailure},
		{"u ta1", []int{}, [2]int{}, parseFailure},
		{"ua1", []int{}, [2]int{}, parseFailure},
		{"uta", []int{}, [2]int{}, parseFailure},
		{"u t1-t1", []int{}, [2]int{}, parseFailure},
		{"u t0-t1", []int{}, [2]int{}, parseFailure},
		{"u t2-t1", []int{}, [2]int{}, parseFailure},
		{"u t1-t" + fmt.Sprint(timeSlicesDisplayed+1), []int{}, [2]int{}, parseFailure},
		// success cases
		{"u t1", []int{1}, [2]int{}, parseSuccess},
		{"ut1", []int{1}, [2]int{}, parseSuccess},
		{"u t3, t6", []int{3, 6}, [2]int{}, parseSuccess},
		{"ut3, t6", []int{3, 6}, [2]int{}, parseSuccess},
		{"u t3,t6", []int{3, 6}, [2]int{}, parseSuccess},
		{"ut3,t6", []int{3, 6}, [2]int{}, parseSuccess},
		{"u t3 t6", []int{3, 6}, [2]int{}, parseSuccess},
		{"ut3 t6", []int{3, 6}, [2]int{}, parseSuccess},
		{"u t3t6", []int{3, 6}, [2]int{}, parseSuccess},
		{"ut3t6", []int{3, 6}, [2]int{}, parseSuccess},
		{"u t7-t10", []int{}, [2]int{7, 10}, parseSuccess},
		{"ut7-t10", []int{}, [2]int{7, 10}, parseSuccess},
		{"u t7-10", []int{}, [2]int{7, 10}, parseSuccess},
		{"ut7-10", []int{}, [2]int{7, 10}, parseSuccess}}
	t.Log("Test: parsing user time unassignment...")
	initRegExp()
	for i, testCase := range testCases {
		timeSlices, timeRange, err := parseUnassignment(testCase.timeUnassignment)
		if !err == testCase.err {
			t.Errorf("Test: parse unassignment FAIL -  parse outcome")
		} else if !reflect.DeepEqual(timeSlices, testCase.timeSlices) {
			t.Errorf("Test: parse unassignment FAIL - time slices in test case %d", i+1)
		} else if timeRange != testCase.timeRange {
			t.Errorf("Test: parse unassignment FAIL - time range in test case %d", i+1)
		} else {
			t.Log("Test: success for unassignment test case " + fmt.Sprint(i+1))
		}
	}
}
