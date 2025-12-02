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
	moves := strings.Split(input, "\n")

	position := 50
	zeroesHit := 0

	for _, inputStr := range moves {
		// wow this is gonna be a lot like that recent everybody codes
		// note: 1 == right and -1 == left
		direction := 1
		if strings.HasPrefix(inputStr, "L") {
			direction = -1
		}
		magnitude, err := strconv.Atoi(strings.TrimLeft(inputStr, "LR"))
		if err != nil {
			panic("bad input")
		}

		position += direction * magnitude
		position %= 100

		if position == 0 {
			zeroesHit++
		}
	}

	fmt.Println("Part 1:", zeroesHit)

}

// thinking - max move
// start at 0, 99 left or 99 right, no trigger
// start at 1, 0 left or 98 right, no trigger
// start at 2, 1 left or 97 right, no trigger
// start at 3, 2 left or 96 right, no trigger
// start at 4, 3 left or 95 right, no trigger
// start at 5, 4 left or 94 right, no trigger
// ...
// start at 95, 94 left or 4 right, no trigger
// start at 96, 95 left or 3 right, no trigger
// start at 97, 96 left or 2 right, no trigger
// start at 98, 97 left or 1 right, no trigger
// start at 99, 98 left or 0 right, no trigger
// all moves outside of this trigger

// at pos 0
// (99-0)%100 == 99 // RIGHT
// (99-0)%100 == 99 // LEFT

// // at pos 1
// (99-1)%100 == 98       // RIGHT
// (99-(100-99))%100 == 0 // LEFT

// // at pos 2
// (99-2)%100 == 97       // RIGHT
// (99-(100-98))%100 == 1 // LEFT

// // at pos 3
// (99-3)%100 == 96       // RIGHT
// (99-(100-97))%100 == 2 // LEFT

// // at pos 4
// (99-4)%100 == 95       // RIGHT
// (99-(100-96))%100 == 3 // LEFT

// // at pos =5
// (99-5)%100 == 94       // RIGHT
// (99-(100-95))%100 == 4 // LEFT

// // ...

// // at pos = 98
// (99-98)%100 == 1        // RIGHT
// (99-(100-98))%100 == 97 //LEFT

// // at pos = 99
// (99-99)%100 == 0        // RIGHT
// (99-(100-99))%100 == 98 //LEFT

// so formula is
// (99 - (100 - currentPosition)) % 100

func part2() {
	moves := strings.Split(input, "\n")

	currentPosition := 50

	zeroesPassed := 0

	for _, inputStr := range moves {
		// fmt.Println("\n--- iter ---")
		direction := 1
		if strings.HasPrefix(inputStr, "L") {
			direction = -1
		}
		rawMagnitude, err := strconv.Atoi(strings.TrimLeft(inputStr, "LR"))
		if err != nil {
			panic("bad input")
		}

		// any move that is greater than 100 will always complete a full spin for each hundred
		magnitude := rawMagnitude
		for magnitude > 100 {
			// fmt.Printf("there is %d left to turn, so clicking for the next turn\n", magnitude)
			magnitude -= 100
			zeroesPassed++
			// fmt.Println("zeroesPassed++, now at", zeroesPassed)
		}

		// if we are just dealing with full turns (multiples of 100)
		// then we have already counted the click to click
		if magnitude == 0 {
			continue
		}

		simplifiedMove := magnitude * direction

		// fmt.Println("currentPosition:", currentPosition)
		// fmt.Println("simplifiedMove:", simplifiedMove)

		maxRightMove := (99 - currentPosition) % 100

		var maxLeftMove int
		if currentPosition == 0 {
			// special case handling for this way of representing the cycle, comes up if running the example too
			maxLeftMove = -99
		} else {
			maxLeftMove = -1 * ((99 - (100 - currentPosition)) % 100)
		}
		// fmt.Println("maxRightMove:", maxRightMove)
		// fmt.Println("maxLeftMove:", maxLeftMove)

		if simplifiedMove < maxLeftMove || simplifiedMove > maxRightMove {
			zeroesPassed++
			// fmt.Println("*** zeroesPassed++, now at", zeroesPassed)
		}

		currentPosition += simplifiedMove
		currentPosition %= 100
		if currentPosition < 0 {
			currentPosition += 100
		}

		// fmt.Println("adjusted position at end:", currentPosition)

	}

	fmt.Println("Part 2:", zeroesPassed)
}

func main() {
	part1()
	part2()
}
