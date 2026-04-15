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
	add           = "+"
	multiply      = "*"
	leftAlligned  = "L"
	rightAlligned = "R"
)

type columnInfo struct {
	operation  string // add or multiply
	allignment string // leftAlligned or rightAlligned
}

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
	lines := strings.Split(input, "\n")
	operationRowIndex := len(lines) - 1

	noNumberRows := len(lines) - 1

	colsInfo := []columnInfo{}

	// build a similar matrix as part 1, but this time keep the numbers as strings!
	// it is easier to check length/significant figures this way. We will do the conversion
	// to int later when we are doing our summing/multiplying.
	numbersStrMatrix := [][]string{}
	for i := range noNumberRows {
		numbersStrRow := strings.Fields(lines[i])
		numbersStrMatrix = append(numbersStrMatrix, numbersStrRow)
	}

	// instead of doing the `operations := strings.Fields(lines[operationRowIndex])` as we do in part 1 we will instead iterate
	// over all the characters in the operations row, since we will use the position of the character to then determine if this
	// set of numbers are left alligned or right alligned. We can tell which orientation the allignment is by looking "up" the column
	// to see if there is any spaces. No spaces means the number is left alligned. Spaces mean right alligned. Edge case is when every
	// number has all the digits, we will see that as being left alligned, but it doesn't matter as in that case there is no additional 0s
	// to add to these numbers.
	// If this concept looks wedged in here, it's because it is. I realised at the end that i didn't account for left vs right allignment!
	for i := range len(lines[operationRowIndex]) {
		maybeOperation := string(lines[operationRowIndex][i])
		if maybeOperation != " " {
			allignment := leftAlligned
			for j := range noNumberRows {
				if string(lines[j][i]) == " " {
					allignment = rightAlligned
					break
				}
			}
			colInfo := columnInfo{
				operation:  maybeOperation,
				allignment: allignment,
			}
			colsInfo = append(colsInfo, colInfo)
		}
	}

	totalTally := 0

	// since we are going right to left, set it up as such in our logic
	// to aid debugging, even though it could be done either way.
	columns := len(numbersStrMatrix[0]) - 1
	for column := columns; column >= 0; column-- {
		// for each column, determine the max figures in a number. This will be used to  determine
		// how many loops we will need to run over each column to perform all the calculations
		maxFigures := 0
		for i := range noNumberRows {
			figures := len(numbersStrMatrix[i][column])
			if figures > maxFigures {
				maxFigures = figures
			}
		}

		numbersForCalculation := []int{}

		if colsInfo[column].allignment == leftAlligned {
			for figure := maxFigures; figure > 0; figure-- {
				digits := ""
				for row := range noNumberRows {
					if len(numbersStrMatrix[row][column]) >= figure {
						digits += string(numbersStrMatrix[row][column][figure-1])
					}
				}
				digitsAsInt, err := strconv.Atoi(digits)
				if err != nil {
					panic("bad input handling!")
				}
				numbersForCalculation = append(numbersForCalculation, digitsAsInt)
			}
		} else { // right alligned
			for figure := 0; figure < maxFigures; figure++ {
				digits := ""
				for row := range noNumberRows {
					if len(numbersStrMatrix[row][column]) > figure {
						digits += string(numbersStrMatrix[row][column][len(numbersStrMatrix[row][column])-1-figure])
					}
				}
				digitsAsInt, err := strconv.Atoi(digits)
				if err != nil {
					panic("bad input handling!")
				}
				numbersForCalculation = append(numbersForCalculation, digitsAsInt)
			}
		}

		// now we can math just like before! Always being careful with the first number.
		// that first number case is why we construct the numbersForCalculation rather than doing it all
		// together above
		tally := numbersForCalculation[0]
		for i := 1; i < len(numbersForCalculation); i++ {
			tally = applyOperation(colsInfo[column].operation, numbersForCalculation[i], tally)
		}
		totalTally += tally
	}

	fmt.Println("part 2:", totalTally)

}

func main() {
	part1()
	part2()
}
