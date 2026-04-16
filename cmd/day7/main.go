package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

type location struct {
	column int
	row    int
}

const (
	startLocation = "S"
	laser         = "|"
	splitter      = "^"
)

func part1() {
	lines := strings.Split(input, "\n")

	rows := len(lines)
	columns := len(lines[0])

	manifold := make(map[location]string)

	// record the starting location and the position of splitters.
	// Everything else at this stage is empty space.
	for y, line := range lines {
		for x, iterRune := range line {
			char := string(iterRune)
			// replace the starting location with a laser...
			if char == startLocation {
				manifold[location{column: x, row: y}] = laser
				continue
			}

			if char == splitter {
				manifold[location{column: x, row: y}] = splitter
			}
		}
	}

	splitCount := 0
	// go through each row, advance the lasers. At the end, count the number of lasers in the last row
	for row := range rows {
		for column := range columns {
			currentLocation := location{column: column, row: row}
			val, _ := manifold[currentLocation]
			if val == laser {
				splitCount = advanceLaser(manifold, currentLocation, splitCount)
			}
		}
	}

	printManifold(manifold, columns, rows)

	fmt.Println("part 1:", splitCount)
}

func advanceLaser(manifold map[location]string, from location, splitCount int) int {
	// laser from location can either go down into blank space, or hit a splitter
	// and split itself into 2 spaces
	down := location{
		column: from.column,
		row:    from.row + 1,
	}

	val, _ := manifold[down]

	// NOTE: checking the input, it is not possible that 2 lasers are right next to each other :)
	// Also note we dont care if there is already a laser in the spot already (probably we will care come part 2)
	if val != splitter {
		manifold[down] = laser
	} else {
		downLeft := location{
			column: from.column - 1,
			row:    from.row + 1,
		}
		downRight := location{
			column: from.column + 1,
			row:    from.row + 1,
		}
		manifold[downLeft] = laser
		manifold[downRight] = laser
		splitCount++
	}

	return splitCount
}

func printManifold(manifold map[location]string, columnLen int, rowLen int) {
	for row := range rowLen {
		rowPretty := ""
		for column := range columnLen {
			val, ok := manifold[location{row: row, column: column}]
			if !ok {
				rowPretty += "."
			} else {
				rowPretty += val
			}
		}
		fmt.Println(rowPretty)
	}
}

func part2() {
}

func main() {
	part1()
	part2()
}
