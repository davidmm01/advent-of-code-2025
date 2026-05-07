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

type tile struct {
	ttype    string // laser or splitter
	lazerNum int    // only applies to lasers
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

	// printManifold(manifold, columns, rows)

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

func advanceLaserPart2(manifold map[location]tile, from location) {
	current, ok := manifold[from]
	if !ok || current.ttype != laser {
		panic("only should be advancing lasers!")
	}

	// laser from location can either go down into blank space, or hit a splitter
	// and split itself into 2 spaces
	down := location{
		column: from.column,
		row:    from.row + 1,
	}

	val, ok := manifold[down]

	if !ok { // nothing there, propagate the magnitude of current laser down
		manifold[down] = tile{ttype: laser, lazerNum: current.lazerNum}
	} else if ok && val.ttype == laser { // theres already a laser there, add the number of lasers
		val.lazerNum += current.lazerNum
		manifold[down] = val
	} else if ok && val.ttype == splitter {
		downLeft := location{
			column: from.column - 1,
			row:    from.row + 1,
		}
		downRight := location{
			column: from.column + 1,
			row:    from.row + 1,
		}

		downLeftTile, leftOk := manifold[downLeft]
		downRightTile, rightOk := manifold[downRight]

		if !leftOk {
			manifold[downLeft] = tile{ttype: laser, lazerNum: current.lazerNum}
		} else {
			downLeftTile.lazerNum += current.lazerNum
			manifold[downLeft] = downLeftTile
		}

		if !rightOk {
			manifold[downRight] = tile{ttype: laser, lazerNum: current.lazerNum}
		} else {
			downRightTile.lazerNum += current.lazerNum
			manifold[downRight] = downRightTile
		}
	}
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

func printManifoldPart2(manifold map[location]tile, columnLen int, rowLen int) {
	for row := range rowLen {
		rowPretty := ""
		for column := range columnLen {
			val, ok := manifold[location{row: row, column: column}]
			// make each pretty print take 2 squares so we can fit bigger numbers stacking,
			// this is sufficient for the sample input:
			// . . . . . . . 1 . . . . . . .
			// . . . . . . . 1 . . . . . . .
			// . . . . . . 1 ^ 1 . . . . . .
			// . . . . . . 1 . 1 . . . . . .
			// . . . . . 1 ^ 2 ^ 1 . . . . .
			// . . . . . 1 . 2 . 1 . . . . .
			// . . . . 1 ^ 3 ^ 3 ^ 1 . . . .
			// . . . . 1 . 3 . 3 . 1 . . . .
			// . . . 1 ^ 4 ^ 3 3 1 ^ 1 . . .
			// . . . 1 . 4 . 3 3 1 . 1 . . .
			// . . 1 ^ 5 ^ 4 3 4 ^ 2 ^ 1 . .
			// . . 1 . 5 . 4 3 4 . 2 . 1 . .
			// . 1 ^ 1 5 4 ^ 7 4 . 2 1 ^ 1 .
			// . 1 . 1 5 4 . 7 4 . 2 1 . 1 .
			// 1 ^ 2 ^10 ^11 ^11 ^ 2 1 1 ^ 1
			// 1 . 2 .10 .11 .11 . 2 1 1 . 1

			if !ok {
				rowPretty += " ."
			} else if val.ttype == splitter {
				rowPretty += " ^"
			} else {
				rowPretty += fmt.Sprintf("%2d", val.lazerNum)
			}
		}
		fmt.Println(rowPretty)
	}
}

func part2() {
	lines := strings.Split(input, "\n")

	rows := len(lines)
	columns := len(lines[0])

	// key diff from part 1: `tile` as map value, rather than just string, since now we also want to track number of lasers living on that tile
	manifold := make(map[location]tile)

	// record the starting location and the position of splitters.
	// Everything else at this stage is empty space.
	for y, line := range lines {
		for x, iterRune := range line {
			char := string(iterRune)
			// replace the starting location with a laser...
			if char == startLocation {
				manifold[location{column: x, row: y}] = tile{ttype: laser, lazerNum: 1}
				continue
			}

			if char == splitter {
				manifold[location{column: x, row: y}] = tile{ttype: splitter}
			}
		}
	}

	// go through each row, advance the lasers
	for row := range rows {
		for column := range columns {
			currentLocation := location{column: column, row: row}
			val, _ := manifold[currentLocation]
			if val.ttype == laser {
				advanceLaserPart2(manifold, currentLocation)
			}
		}
	}

	// count up the number of paths at the end bum summing up all the lasers we see at the end
	totalLasers := 0
	for col := range columns {
		loc := location{column: col, row: rows - 1}
		countTile, ok := manifold[loc]
		if ok && countTile.ttype == laser {
			totalLasers += countTile.lazerNum
		}
	}

	// printManifoldPart2(manifold, columns, rows)

	fmt.Println("part 2:", totalLasers)
}

// TODO: make generic! lots of opportunity to refactor and reduce the use of part1/part2 functions for this day

func main() {
	part1()
	part2()
}
