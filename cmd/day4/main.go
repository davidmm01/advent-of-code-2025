package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"strings"
)

//go:embed input.txt
var input string

const (
	emptySpace = "."
	roll       = "@"
)

type mapConstraints struct {
	xMin int
	xMax int
	yMin int
	yMax int
}

type coordinate struct {
	x int
	y int
}

func part1() {
	lines := strings.Split(input, "\n")

	constraints := mapConstraints{
		xMin: 0,
		xMax: len(lines[0]) - 1,
		yMin: 0,
		yMax: len(lines) - 1,
	}

	numCanAccessRolls := 0

	for y, line := range strings.Split(input, "\n") {
		for x, charRune := range line {
			if string(charRune) != roll {
				continue
			}

			coords := getCoordinatesToCheck(coordinate{x, y})
			emptySpots := 0
			for _, coord := range coords {
				if isCoordOutOfBounds(coord, constraints) || string(lines[coord.y][coord.x]) == emptySpace {
					emptySpots++
				}
			}

			if emptySpots >= 5 {
				numCanAccessRolls++
			}
		}
	}

	fmt.Println("Part 1:", numCanAccessRolls)

}

func isCoordOutOfBounds(in coordinate, constraints mapConstraints) bool {
	return in.x < constraints.xMin ||
		in.y < constraints.yMin ||
		in.x > constraints.xMax ||
		in.y > constraints.yMax
}

func getCoordinatesToCheck(in coordinate) []coordinate {
	up := coordinate{in.x, in.y - 1}
	down := coordinate{in.x, in.y + 1}
	left := coordinate{in.x - 1, in.y}
	right := coordinate{in.x + 1, in.y}
	upLeft := coordinate{in.x - 1, in.y - 1}
	upRight := coordinate{in.x + 1, in.y - 1}
	downLeft := coordinate{in.x - 1, in.y + 1}
	downRight := coordinate{in.x + 1, in.y + 1}
	return []coordinate{up, down, left, right, upLeft, upRight, downLeft, downRight}
}

func part2() {
	lines := strings.Split(input, "\n")

	constraints := mapConstraints{
		xMin: 0,
		xMax: len(lines[0]) - 1,
		yMin: 0,
		yMax: len(lines) - 1,
	}

	numCanAccessRolls := 0

	canAccessCoords := []coordinate{}

	for { // infinite loop, terminate by condition at end
		for y, line := range lines {
			for x, charRune := range line {
				if string(charRune) != roll {
					continue
				}

				coords := getCoordinatesToCheck(coordinate{x, y})
				emptySpots := 0
				for _, coord := range coords {
					if isCoordOutOfBounds(coord, constraints) || string(lines[coord.y][coord.x]) == emptySpace {
						emptySpots++
					}
				}

				if emptySpots >= 5 {
					numCanAccessRolls++
					canAccessCoords = append(canAccessCoords, coordinate{x, y})
				}
			}
		}

		if len(canAccessCoords) == 0 {
			break
		}

		for _, coord := range canAccessCoords {
			// i hate the way with all these grid questions i end up with indexes of [y][x] instead of [x][y]
			lines[coord.y] = lines[coord.y][:coord.x] + emptySpace + lines[coord.y][coord.x+1:]
		}

		canAccessCoords = nil
	}

	fmt.Println("part 2:", numCanAccessRolls)
}

func main() {
	part1()
	part2()
}
