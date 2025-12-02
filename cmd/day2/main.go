package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type LeRange struct {
	Lower int
	Upper int
}

func part2() {
	ranges := []LeRange{}

	// wrangle input
	for _, line := range strings.Split(input, "\n") {
		for _, aRangeStr := range strings.Split(line, ",") {
			aRangeStrSlice := strings.Split(aRangeStr, "-")
			lower, err := strconv.Atoi(aRangeStrSlice[0])
			if err != nil {
				panic("death")
			}
			upper, err := strconv.Atoi(aRangeStrSlice[1])
			if err != nil {
				panic("death")
			}
			ranges = append(ranges, LeRange{Lower: lower, Upper: upper})
		}
	}

	sum := 0
	for _, aRange := range ranges {
		violators := getViolators(aRange)
		for _, violator := range violators {
			sum += violator
		}

	}

	fmt.Println("part 2:", sum)
}

func getViolators(aRange LeRange) []int {
	violators := []int{}
	for i := aRange.Lower; i <= aRange.Upper; i++ {
		// observation: numbers must be of even number of digits, else they are not in the running
		// Special case, can be an odd number only if every digit is the same!
		// for each number, if it has an even length, break it in half and check the two
		// then break the OG number in half again, and check the four
		// keep breaking OG number
		if isViolator(i) {
			violators = append(violators, i)
		}

	}
	// fmt.Println("violators", violators)
	return violators
}

// this is how i initially did it, but im not meant to count numbers like 999 so its wrong!
// but this might be needed for part 2 i suspect
// omg and i was right lol
// i need to read first before trying to solve.
func isViolator(input int) bool {
	inputStr := strconv.Itoa(input)

	// every number (even odd ones) might satisfy that each number is the same!
	numSegments := len(inputStr)

	// we could optomise here by only checking half of the numbers
	for i := numSegments; i > 1; i-- {
		if numSegments%i == 0 {
			chunks := makeChunks(inputStr, i)
			if checkChunk(chunks) {
				return true
			}
		}
	}

	return false
}

func checkChunk(chunks []string) bool {
	for i := 0; i < len(chunks)-1; i++ {
		if chunks[i] != chunks[i+1] {
			return false
		}
	}
	return true
}

func makeChunks(inputNumber string, segments int) []string {
	inputNumberLen := len(inputNumber)

	if inputNumberLen%segments != 0 {
		panic(fmt.Sprintf("you shouldn't call me with inputNumber=%s and segments=%d", inputNumber, segments))
	}

	chunkSize := inputNumberLen / segments

	chunks := []string{}
	for i := 0; i < inputNumberLen; i += chunkSize {
		chunks = append(chunks, inputNumber[i:i+chunkSize])
	}

	return chunks
}

func part1() {
	ranges := []LeRange{}

	// wrangle input
	for _, line := range strings.Split(input, "\n") {
		for _, aRangeStr := range strings.Split(line, ",") {
			aRangeStrSlice := strings.Split(aRangeStr, "-")
			lower, err := strconv.Atoi(aRangeStrSlice[0])
			if err != nil {
				panic("death")
			}
			upper, err := strconv.Atoi(aRangeStrSlice[1])
			if err != nil {
				panic("death")
			}
			ranges = append(ranges, LeRange{Lower: lower, Upper: upper})
		}
	}

	// fmt.Println("ranges found:", ranges)

	sum := 0
	for _, aRange := range ranges {
		hits := checkRange(aRange)
		// fmt.Println("hits:", hits)
		for _, hit := range hits {
			sum += hit
		}

	}

	fmt.Println("part 1:", sum)
}

func checkRange(aRange LeRange) []int {
	hits := []int{}

	for current := aRange.Lower; current <= aRange.Upper; current++ {
		// fmt.Println("checking current:", current)
		currentAsString := strconv.Itoa(current)
		theLen := len(currentAsString)

		// if the length is even, then we can just compare the first number of each half, then the second number, etc...
		if theLen%2 == 0 {
			allMatch := true
			for j := 0; j < theLen/2; j++ {
				// fmt.Printf("checking %s with %s, its %v\n", string(currentAsString[j]), string(currentAsString[j+(theLen/2)]), string(currentAsString[j]) != string(currentAsString[j+(theLen/2)]))
				if string(currentAsString[j]) != string(currentAsString[j+(theLen/2)]) {
					allMatch = false
					break
				}

			}
			// fmt.Println("im matching all!")
			if allMatch {
				hits = append(hits, current)
			}
		}

	}
	return hits
}

func main() {
	part1()
	part2()
}
