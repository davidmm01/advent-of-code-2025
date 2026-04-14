package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

const (
	add      = "+"
	multiply = "*"
)

func part1() {
	lines := strings.Split(input, "\n")
	noNumberRows := len(lines) - 1
	operationRowIndex := len(lines) - 1

	// build up a matrix of the numbers to work on, so we can easily reference columns,
	// i.e. can do numbersMatrix[row][column]
	numbersMatrix := [][]int{}
	for i := range noNumberRows {
		numbersAsStr := strings.Fields(lines[i])
		numberRow := []int{}
		for _, numberAsStr := range numbersAsStr {
			numberAsInt, err := strconv.Atoi(numberAsStr)
			if err != nil {
				panic(fmt.Sprintf("bad input parsing logic! '%s' is not an int", numberAsStr))
			}
			numberRow = append(numberRow, numberAsInt)
		}
		numbersMatrix = append(numbersMatrix, numberRow)
	}

	// iterate over operations and do them
	totalTally := 0
	operations := strings.Fields(lines[operationRowIndex])
	for column, operation := range operations {
		// tally needs to start at the first number rather than 0 so we dont get caught in a mult*0 loop
		tally := numbersMatrix[0][column]
		for row := 1; row < noNumberRows; row++ {
			tally = applyOperation(operation, numbersMatrix[row][column], tally)
		}
		totalTally += tally
	}

	fmt.Println("part 1:", totalTally)
}

func applyOperation(operation string, incoming, tally int) int {
	if operation == add {
		return incoming + tally
	}
	return incoming * tally
}

func part2() {
}

func main() {
	part1()
	part2()
}
