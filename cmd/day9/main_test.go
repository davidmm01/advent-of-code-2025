package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_squareSize(t *testing.T) {
	tests := []struct {
		loc1     location
		loc2     location
		expected int
	}{
		{
			loc1:     location{column: 1, row: 1},
			loc2:     location{column: 1, row: 1},
			expected: 1,
		},
		{
			loc1:     location{column: 1, row: 1},
			loc2:     location{column: 2, row: 2},
			expected: 4,
		},
		{
			loc1:     location{column: 1, row: 1},
			loc2:     location{column: 2, row: 1},
			expected: 2,
		},
		{
			loc1:     location{column: 2, row: 5},
			loc2:     location{column: 9, row: 7},
			expected: 24,
		},
		{
			loc1:     location{column: 7, row: 3},
			loc2:     location{column: 2, row: 3},
			expected: 6,
		},
		{
			loc1:     location{column: 11, row: 1},
			loc2:     location{column: 2, row: 5},
			expected: 50,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equal(t, tt.expected, squareSize(tt.loc1, tt.loc2))
			assert.Equal(t, tt.expected, squareSize(tt.loc2, tt.loc1)) // check order doesn't matter
		})
	}
}

func Test_getColumnMinMax(t *testing.T) {
	row := 55 // this should always be the same for the calling set

	tests := []struct {
		locs           []location
		expectedMinCol location
		expectedMaxCol location
	}{
		{
			locs: []location{
				{column: 12, row: row},
				{column: 14, row: row},
				{column: 20, row: row},
				{column: 400, row: row},
			},
			expectedMinCol: location{column: 12, row: row},
			expectedMaxCol: location{column: 400, row: row},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actualMinCol, actualMaxCol := getColumnMinMax(tt.locs)
			assert.Equal(t, tt.expectedMinCol, actualMinCol)
			assert.Equal(t, tt.expectedMaxCol, actualMaxCol)
		})
	}
}

func Test_getRowMinMax(t *testing.T) {
	column := 55 // this should always be the same for the calling set

	tests := []struct {
		locs           []location
		expectedMinRow location
		expectedMaxRow location
	}{
		{
			locs: []location{
				{column: column, row: 54},
				{column: column, row: 22},
				{column: column, row: 786},
				{column: column, row: 454},
			},
			expectedMinRow: location{column: column, row: 22},
			expectedMaxRow: location{column: column, row: 786},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actualMinRow, actualMaxRow := getRowMinMax(tt.locs)
			assert.Equal(t, tt.expectedMinRow, actualMinRow)
			assert.Equal(t, tt.expectedMaxRow, actualMaxRow)
		})
	}
}

func Test_getOtherCornersOfSquare(t *testing.T) {
	tests := []struct {
		loc1               location
		loc2               location
		expectedNewCorner1 location
		expectedNewCorner2 location
	}{
		{
			loc1:               location{row: 9, column: 5},
			loc2:               location{row: 11, column: 7},
			expectedNewCorner1: location{row: 9, column: 7},
			expectedNewCorner2: location{row: 11, column: 5},
		},
		{
			loc1:               location{row: 1, column: 1},
			loc2:               location{row: 2, column: 2},
			expectedNewCorner1: location{row: 1, column: 2},
			expectedNewCorner2: location{row: 2, column: 1},
		},
		{
			loc1:               location{row: 1, column: 9},
			loc2:               location{row: 1, column: 10},
			expectedNewCorner1: location{row: 1, column: 10},
			expectedNewCorner2: location{row: 1, column: 9},
		},
		{
			loc1:               location{row: 1, column: 8},
			loc2:               location{row: 2, column: 8},
			expectedNewCorner1: location{row: 1, column: 8},
			expectedNewCorner2: location{row: 2, column: 8},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			actualNewCorner1, actualNewCorner2 := getOtherCornersOfSquare(tt.loc1, tt.loc2)
			assert.Equal(t, tt.expectedNewCorner1, actualNewCorner1)
			assert.Equal(t, tt.expectedNewCorner2, actualNewCorner2)
		})
	}
}
