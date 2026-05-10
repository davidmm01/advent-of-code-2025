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
			loc1:     location{x: 1, y: 1},
			loc2:     location{x: 1, y: 1},
			expected: 1,
		},
		{
			loc1:     location{x: 1, y: 1},
			loc2:     location{x: 2, y: 2},
			expected: 4,
		},
		{
			loc1:     location{x: 1, y: 1},
			loc2:     location{x: 2, y: 1},
			expected: 2,
		},
		{
			loc1:     location{x: 2, y: 5},
			loc2:     location{x: 9, y: 7},
			expected: 24,
		},
		{
			loc1:     location{x: 7, y: 3},
			loc2:     location{x: 2, y: 3},
			expected: 6,
		},
		{
			loc1:     location{x: 11, y: 1},
			loc2:     location{x: 2, y: 5},
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
