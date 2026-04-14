package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const endSeq = "\n\n1" // we dont are about this part much, just to make a test data a bit easier to read.

func Test_Part2_Main(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedSum    int
		expectedRanges []ARange
	}{
		{
			name:           "sample input provided",
			input:          "3-5\n10-14\n16-20\n12-18" + endSeq,
			expectedSum:    14,
			expectedRanges: []ARange{{3, 5}, {10, 20}},
		},
		{
			name:           "absorb, expanding the lower bound only",
			input:          "3-5\n2-4" + endSeq,
			expectedSum:    4,
			expectedRanges: []ARange{{2, 5}},
		},
		{
			name:           "absorb, expanding the lower bound only with a single element slice",
			input:          "3-5\n2-2" + endSeq,
			expectedSum:    4,
			expectedRanges: []ARange{{2, 2}, {3, 5}}, // yes, this looks weird but its fine because the sum is still correct
		},
		{
			name:           "absorb, expanding the upper bound only",
			input:          "3-5\n4-6" + endSeq,
			expectedSum:    4,
			expectedRanges: []ARange{{3, 6}},
		},
		{
			name:           "absorb, expanding both bounds",
			input:          "3-5\n2-6" + endSeq,
			expectedSum:    5,
			expectedRanges: []ARange{{2, 6}},
		},
		{
			name:           "absorb, no expansion",
			input:          "3-6\n4-5" + endSeq,
			expectedSum:    4,
			expectedRanges: []ARange{{3, 6}},
		},
		{
			name:           "absorb, no expansion, single element slice",
			input:          "3-5\n4-4" + endSeq,
			expectedSum:    3,
			expectedRanges: []ARange{{3, 5}},
		},
		{
			name:           "absorb multiple data points at once",
			input:          "1-2\n3-5\n7-8\n9-10\n12-15\n3-12" + endSeq,
			expectedSum:    15,
			expectedRanges: []ARange{{1, 2}, {3, 15}},
		},
		{
			name:           "absorb multiple data points at once",
			input:          "1-2\n3-4\n5-6\n7-8" + endSeq,
			expectedSum:    8,
			expectedRanges: []ARange{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
		},
		{
			name:           "repetition doesnt matter",
			input:          "1-2\n3-4\n5-6\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8\n7-8" + endSeq,
			expectedSum:    8,
			expectedRanges: []ARange{{1, 2}, {3, 4}, {5, 6}, {7, 8}},
		},
		{
			name:           "absorb all data points, then add a new low",
			input:          "5-10\n15-20\n25-30\n5-29\n1-2" + endSeq,
			expectedSum:    28,
			expectedRanges: []ARange{{1, 2}, {5, 30}},
		},
		{
			name:           "absorb all data points, then add a new high",
			input:          "5-10\n15-20\n25-30\n5-29\n40-50" + endSeq,
			expectedSum:    37,
			expectedRanges: []ARange{{5, 30}, {40, 50}},
		},
		// sadly i got stuck and couldnt figure out why i was failing, i couldn't think of more edge cases and wasted hours.
		// Asked claude, this was the missing edge case:
		//
		// BUG: getIndexInfo computed the wrong insertIndex when the incoming range belonged at
		// position 0 and there were already 2+ ranges in the list.
		// ^ davo note: super annoying, probably would not have found this myself without just trying heaps of random combinations... really sad now lol.
		//
		// The old logic had two cases:
		//   - generic case (not last element): set insertIndex if incomingRange.Lower fell between
		//     ranges[i].Lower and ranges[i+1].Lower
		//   - last element case: set insertIndex = i if incomingRange.Lower < ranges[i].Lower, else i+1
		//
		// The flaw: when incomingRange.Lower was less than ALL existing ranges, the generic case
		// never fired (e.g. 5 < 1 is false). Then the last-element handler set insertIndex = lastIndex
		// (e.g. 1 for a 2-element slice) instead of 0. The range was inserted out of order, which
		// corrupted absorbRanges on subsequent steps — it would absorb the wrong pair of ranges and
		// silently drop coverage.
		//
		// The fix: track insertIndex as "first index where ranges[i].Lower >= incomingRange.Lower",
		// defaulting to len(ranges) (append to end). One simple condition handles all cases.
		//
		// now back to me: ALSO LET ME SAY IN MY SHAME that none of the question was given to claude. I didnt even give him a real prompt, i just gave him
		// my code and said "fix it" so i was close enough that it could see the bug lol.
		{
			name:           "insert at front when 2+ ranges already exist, then merge",
			input:          "5-10\n15-20\n1-2\n1-6" + endSeq,
			expectedSum:    16,
			expectedRanges: []ARange{{1, 10}, {15, 20}},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test case: %s", tt.name), func(t *testing.T) {
			acutalSum, actualRanges := part2(tt.input)
			assert.Equal(t, tt.expectedSum, acutalSum)
			assert.Equal(t, tt.expectedRanges, actualRanges)
		})
	}
}

func Test_hasOverlap(t *testing.T) {
	tests := []struct {
		name          string
		checkRange    ARange
		incomingRange ARange
		expected      bool
	}{
		{
			name:          "overlap on min side",
			checkRange:    ARange{3, 5},
			incomingRange: ARange{2, 4},
			expected:      true,
		},
		{
			name:          "overlap on max side",
			checkRange:    ARange{3, 5},
			incomingRange: ARange{4, 6},
			expected:      true,
		},
		{
			name:          "incoming range is fully contained within check range",
			checkRange:    ARange{2, 6},
			incomingRange: ARange{3, 5},
			expected:      true,
		},
		{
			name:          "check range is fully contained within incoming range",
			checkRange:    ARange{3, 5},
			incomingRange: ARange{2, 6},
			expected:      true,
		},
		{
			name:          "fully contained single element slice",
			checkRange:    ARange{2, 6},
			incomingRange: ARange{4, 4},
			expected:      true,
		},
		{
			name:          "reverse of fully contained single element slice",
			checkRange:    ARange{4, 4},
			incomingRange: ARange{2, 6},
			expected:      true,
		},
		{
			name:          "no overlap",
			checkRange:    ARange{2, 6},
			incomingRange: ARange{7, 8},
			expected:      false,
		},
		{
			name:          "no overlap",
			checkRange:    ARange{1, 1},
			incomingRange: ARange{2, 2},
			expected:      false,
		},
		{
			name:          "no overlap",
			checkRange:    ARange{2, 2},
			incomingRange: ARange{1, 1},
			expected:      false,
		},
		{
			name:          "no overlap",
			checkRange:    ARange{1, 2},
			incomingRange: ARange{3, 4},
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("test case: %s", tt.name), func(t *testing.T) {
			assert.Equal(t, tt.expected, hasOverlap(tt.checkRange, tt.incomingRange))
		})
	}
}
