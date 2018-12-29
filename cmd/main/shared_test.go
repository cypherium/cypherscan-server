package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetMissedNumbers is to test getMissedNumbers
func TestGetMissedNumbers(t *testing.T) {
	cases := []struct {
		startWith         int64
		pageSize          int
		numbersAlreadyGot []int64
		expectedResult    []int64
	}{
		{10, 6, []int64{9, 6}, []int64{10, 8, 7, 5}},
		{10, 2, []int64{10}, []int64{9}},
		{6, 3, []int64{4}, []int64{6, 5}},
		{11, 5, []int64{}, []int64{11, 10, 9, 8, 7}},
	}

	for i, c := range cases {
		assert.Equal(t, c.expectedResult, getMissedNumbers(c.startWith, c.pageSize, c.numbersAlreadyGot), fmt.Sprintf("Failed on %d", i))
	}
}
