package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"slices"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type ARange struct {
	Lower int
	Upper int
}

func part1() {
	lines := strings.Split(input, "\n")
	breakSeen := false
	ranges := []ARange{}

	goodIngredientCount := 0

	for _, line := range lines {
		if line == "" {
			breakSeen = true
			continue
		}

		if !breakSeen {
			rangeSlice := strings.Split(line, "-")
			lower, err := strconv.Atoi(rangeSlice[0])
			if err != nil {
				panic(fmt.Sprintf("bad logic handling input - this isnt a range: %s", line))
			}
			upper, err := strconv.Atoi(rangeSlice[1])
			if err != nil {
				panic(fmt.Sprintf("bad logic handling input - this isnt a range: %s", line))
			}
			ranges = append(ranges, ARange{Lower: lower, Upper: upper})
		} else {
			ingredient, err := strconv.Atoi(line)
			if err != nil {
				panic(fmt.Sprintf("bad logic handling input - this isnt a ingredient to check: %s", line))
			}
			for _, aRange := range ranges {
				if aRange.Lower <= ingredient && aRange.Upper >= ingredient {
					goodIngredientCount++
					break
				}
			}
		}
	}

	fmt.Println("part 1:", goodIngredientCount)
}

// exploration on example:
//
// see 3-5, no existing ranges, add new range: [3,5]
// 10-14, check if overlaps with existing range, no overlap ranges, add new range: [3,5], [10,14]
// 16-20, check if overlaps with existing ranges, no overlap ranges, add new range: [3,5], [10,14], [16-20]
// 12-18, check if overlaps with existing ranges, two overlap ranges, blend with all affected ranges at once: [3,5], [10,20]
// algo for blending with all ranges
// relies on the ranges being tracked to be in order, so we will always insert in order

func part2(inny string) (int, []ARange) {
	lines := strings.Split(inny, "\n")
	ranges := []ARange{}

	for _, line := range lines {
		if line == "" {
			break
		}

		// parse the input
		rangeSlice := strings.Split(line, "-")
		lower, err := strconv.Atoi(rangeSlice[0])
		if err != nil {
			panic(fmt.Sprintf("bad logic handling input - this isnt a range: %s", line))
		}
		upper, err := strconv.Atoi(rangeSlice[1])
		if err != nil {
			panic(fmt.Sprintf("bad logic handling input - this isnt a range: %s", line))
		}
		newRange := ARange{Lower: lower, Upper: upper}

		// get overlapping ranges with the newRange
		overlappingIndexes, insertIndex := getIndexInfo(newRange, ranges)

		// replace overlapping indexes within ranges with the new expanded range,
		// otherwise simply insert in order
		if len(overlappingIndexes) > 0 {
			ranges = absorbRanges(newRange, overlappingIndexes, ranges)
		} else {
			ranges = slices.Insert(ranges, insertIndex, newRange)
		}
	}

	sum := 0
	for _, aRange := range ranges {
		sum += aRange.Upper - aRange.Lower + 1 // +1 is because its inclusive, i.e. a range of [5,5] should be size 1 not 0
	}

	return sum, ranges
}

func absorbRanges(incomingRange ARange, overlappingIndexes []int, ranges []ARange) []ARange {
	lowestOverlappingIndex := overlappingIndexes[0]
	highestOverlappingIndex := overlappingIndexes[len(overlappingIndexes)-1]

	// determine new floor and ceiling of the absorbing range
	absorbedLow := min(incomingRange.Lower, ranges[lowestOverlappingIndex].Lower)
	absorbedHigh := max(incomingRange.Upper, ranges[highestOverlappingIndex].Upper)
	newRange := ARange{Lower: absorbedLow, Upper: absorbedHigh}

	// remove redundant ranges
	ranges = append(ranges[:lowestOverlappingIndex], ranges[highestOverlappingIndex+1:]...)

	// insert into the correct spot
	if len(ranges) == 0 {
		ranges = []ARange{newRange}
	} else {
		ranges = slices.Insert(ranges, lowestOverlappingIndex, newRange)
	}
	return ranges
}

func getIndexInfo(incomingRange ARange, ranges []ARange) (overlapIndexes []int, insertIndex int) {
	overlapIndexes = []int{}
	insertIndex = len(ranges) // default: insert at end

	for i, checkRange := range ranges {
		// determine overlap indexes
		if hasOverlap(checkRange, incomingRange) {
			overlapIndexes = append(overlapIndexes, i)
		}

		// find the first position where the existing range's lower >= incomingRange's lower
		if checkRange.Lower >= incomingRange.Lower && i < insertIndex {
			insertIndex = i
		}
	}

	return overlapIndexes, insertIndex
}

// this feels like there would deffs be somethin more efficient but im over this now
func hasOverlap(checkRange ARange, incomingRange ARange) bool {
	checkFullyContained := (incomingRange.Lower <= checkRange.Lower) && (incomingRange.Upper >= checkRange.Upper)
	incomingFullyContained := (checkRange.Lower <= incomingRange.Lower) && (checkRange.Upper >= incomingRange.Upper)
	//                 incoming lower is outside of check range      incoming upper is inside check range
	lowerExpansion := (incomingRange.Lower <= checkRange.Lower) && ((incomingRange.Upper <= checkRange.Upper) && (incomingRange.Upper >= checkRange.Lower))
	//                 incoming upper is outside check range         incoming lower is inside check range
	upperExpansion := (incomingRange.Upper >= checkRange.Upper) && ((incomingRange.Lower >= checkRange.Lower) && (incomingRange.Lower <= checkRange.Upper))
	return checkFullyContained || incomingFullyContained || lowerExpansion || upperExpansion
}

func main() {
	part1()
	res, _ := part2(input)
	fmt.Println("part 2:", res)
}
