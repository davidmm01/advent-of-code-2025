package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkPart1(b *testing.B) {
	for range b.N {
		part1()
	}
}

func BenchmarkPart2(b *testing.B) {
	for range b.N {
		part2()
	}
}

func Test_straightLineDistance(t *testing.T) {
	tests := []struct {
		loc1     location
		loc2     location
		expected float64
	}{
		{
			// got these test case from a google search
			loc1: location{1, 2, 3},
			loc2: location{4, 6, 8},
			// ~7.07, so this looks good
			expected: float64(7.0710678118654755),
		},
		{
			// got these test case from a google search
			loc1: location{0, 0, 0},
			loc2: location{1, 1, 1},
			// ~1.732, so this looks good
			expected: float64(1.7320508075688772),
		},
		{
			loc1:     location{431, 825, 988},
			loc2:     location{425, 690, 689},
			expected: float64(328.11888089532425),
		},
		{
			loc1:     location{431, 825, 988},
			loc2:     location{162, 817, 812},
			expected: float64(321.560258738545),
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			assert.Equal(t, tt.expected, straightLineDistance(tt.loc1, tt.loc2))
		})
	}
}
