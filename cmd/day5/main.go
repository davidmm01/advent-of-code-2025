package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
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

func part2() {

}

func main() {
	part1()
	part2()
}
