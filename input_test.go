package main

import (
	"fmt"
	"reflect"
	"testing"
)

// TestParseTimeEntry - test user intput of a time to activity entry
func TestParseTimeEntry(t *testing.T) {

	type testCase struct {
		timeEntry  string
		timeSlices []int
		timeRange  [2]int
		activity   int
		err        bool
	}

	testCases := []testCase{
		// failure cases
		testCase{"t1", []int{}, [2]int{}, 0, true},
		testCase{"t1 a", []int{}, [2]int{}, 0, true},
		testCase{"t1a", []int{}, [2]int{}, 0, true},
		testCase{"t a1", []int{}, [2]int{}, 0, true},
		testCase{"ta1", []int{}, [2]int{}, 0, true},
		testCase{"ta", []int{}, [2]int{}, 0, true},
		testCase{"t1-t1 a1", []int{}, [2]int{}, 0, true},
		testCase{"t0-t1 a1", []int{}, [2]int{}, 0, true},
		testCase{"t2-t1 a1", []int{}, [2]int{}, 0, true},
		testCase{"t1-t" + fmt.Sprint(timeSlicesDisplayed+1) + " a1", []int{}, [2]int{}, 0, true},
		// success cases
		testCase{"t1 a1", []int{1}, [2]int{}, 1, false},
		testCase{"t1a1", []int{1}, [2]int{}, 1, false},
		testCase{"t3, t6 a2", []int{3, 6}, [2]int{}, 2, false},
		testCase{"t3,t6 a2", []int{3, 6}, [2]int{}, 2, false},
		testCase{"t3 t6 a2", []int{3, 6}, [2]int{}, 2, false},
		testCase{"t3t6 a2", []int{3, 6}, [2]int{}, 2, false},
		testCase{"t3t6a2", []int{3, 6}, [2]int{}, 2, false},
		testCase{"t7-t10 a5", []int{}, [2]int{7, 10}, 5, false},
		testCase{"t7-t10a5", []int{}, [2]int{7, 10}, 5, false}}
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
