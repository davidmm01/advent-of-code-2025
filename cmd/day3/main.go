package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

func part1() {
	banks := strings.Split(input, "\n")

	sumOfVoltages := 0

	for _, bank := range banks {
		biggestNumber := 0
		biggestNumberIndex := -1

		secondBiggest := 0
		secondBiggestIndex := -1

		for i := 0; i < len(bank); i++ {
			candidate, err := strconv.Atoi(string(bank[i]))
			if err != nil {
				panic("boom")
			}

			// note: strictly greater, we need lower index to trump higher index for our logic
			if candidate > biggestNumber {
				secondBiggest = biggestNumber
				secondBiggestIndex = biggestNumberIndex
				biggestNumber = candidate
				biggestNumberIndex = i
			} else if candidate > secondBiggest {
				secondBiggest = candidate
				secondBiggestIndex = i
			}
		}

		var voltage int
		var err error
		if biggestNumberIndex < secondBiggestIndex {
			// if the first and second biggest numbers are already ordered, then its trivial
			voltage, err = strconv.Atoi(string(bank[biggestNumberIndex]) + string(bank[secondBiggestIndex]))
			if err != nil {
				panic("oopsies")
			}
		} else if biggestNumberIndex == len(bank)-1 {
			// if the biggest number is the last number, then we already have the answer
			voltage, err = strconv.Atoi(string(bank[secondBiggestIndex]) + string(bank[biggestNumberIndex]))
			if err != nil {
				panic("oopsies")
			}

		} else {
			// since the second largest can't be used as we can't change the ordering,
			// need to look again using the biggest number as reference.
			// we only need to check the numbers after it
			secondBiggestStrictlyAfterBiggest := 0
			for i := biggestNumberIndex + 1; i < len(bank); i++ {
				candidate, err := strconv.Atoi(string(bank[i]))
				if err != nil {
					panic("boom")
				}

				if candidate > secondBiggestStrictlyAfterBiggest {
					secondBiggestStrictlyAfterBiggest = candidate
				}
			}
			voltage, err = strconv.Atoi(fmt.Sprintf("%d%d", biggestNumber, secondBiggestStrictlyAfterBiggest))
			if err != nil {
				panic("fook")
			}
		}

		// fmt.Printf("for bank %s, voltage was %d\n", bank, voltage)

		sumOfVoltages += voltage
	}
	fmt.Println("part 1:", sumOfVoltages)
}

func part2() {
	banks := strings.Split(input, "\n")

	sumOfVoltages := 0
	for _, bank := range banks {
		// fmt.Printf("\n--- bank %s\n", bank)

		iLowerBound := 0
		voltageString := ""

		for voltageMembersRemaining := 11; voltageMembersRemaining > -1; voltageMembersRemaining-- {

			biggestNumber := 0
			biggestNumberIndex := -1

			// fmt.Printf("for digit number %d we search between indexes [%d,%d]\n", 12-voltageMembersRemaining, iLowerBound, len(bank)-voltageMembersRemaining)
			for i := iLowerBound; i < len(bank)-voltageMembersRemaining; i++ {
				candidate, err := strconv.Atoi(string(bank[i]))
				if err != nil {
					panic("boom")
				}

				// note: strictly greater, we need lower index to trump higher index for our logic
				if candidate > biggestNumber {
					// fmt.Printf("cadidate %d is bigger than biggest %d\n", candidate, biggestNumber)
					biggestNumber = candidate
					biggestNumberIndex = i
				} else {
					// fmt.Printf("cadidate %d is smaller than biggest %d\n", candidate, biggestNumber)
				}
			}

			voltageString = voltageString + string(bank[biggestNumberIndex])
			// shift the window in which we can search for the next biggest number
			iLowerBound = biggestNumberIndex + 1

			// fmt.Printf("biggestNumber is %d at index %d, iLowerBound now %d\n", biggestNumber, biggestNumberIndex, iLowerBound)
		}

		// fmt.Println("voltageString:", voltageString)

		voltage, err := strconv.Atoi(voltageString)
		if err != nil {
			panic("boom")
		}

		sumOfVoltages += voltage
	}

	fmt.Println("part 2:", sumOfVoltages)
}

func main() {
	part1()
	part2()
}
