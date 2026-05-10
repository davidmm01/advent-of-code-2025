package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"math"
	"strconv"
	"strings"
)

//go:embed input.txt
var input string

type location struct {
	x int
	y int
}

func part1() {
	// Part 1 thinkings:
	// O(n^2) brute force is trivial (find each square size by comparing each point to every other point)
	// BUT we can probably do it in O(n) by:
	// iterating over each point once and recording...
	// record the points closest to the:
	// - bottom right
	// - top left
	// - top right
	// - bottom left
	// make a square with top left and bottom right, or with top right and bottom left, and see which is the biggest
	// I think this will always guarantee the largest possible square?
	// lets find out...

	// Based on the map they give us, this is what i mean by my labels:
	// top left                top right
	//           ..............
	//           .......#...#..
	//           ..............
	//           ..#....#......
	//           ..............
	//           ..#......#....
	//           ..............
	//           .........#.#..
	//           ..............
	// bottom left             bottom right

	// how to determine which point is most bottom left?

	// ..............
	// ..............
	// ........x.....
	// .x............
	// ..............
	// ..............
	// ....x.........
	// ..............
	// ..............

	// The point at [1, 4] is considered more of a bottom left point than [4,6] by looking at the area of the triangles that get made
	// with it and the corner. Instead we need the straght line distance to the corner via pythagoras? Yes this looks good and works!
	locs := []location{}

	for _, line := range strings.Split(input, "\n") {
		xy := strings.Split(line, ",")
		x, err1 := strconv.Atoi(xy[0])
		y, err2 := strconv.Atoi(xy[1])
		if err1 != nil || err2 != nil {
			panic(fmt.Sprintf("bad input data, line: '%s'", line))
		}
		locs = append(locs, location{x: x, y: y})
	}

	var topLeft location
	topLeftBound := location{x: 0, y: 0}
	topLeftSmallestDist := float64(1000000000000000000)

	var topRight location
	topRightBound := location{x: 100_000, y: 0}
	topRightSmallestDist := float64(1000000000000000000)

	var bottomLeft location
	bottomLeftBound := location{x: 0, y: 100_00}
	bottomLeftSmallestDist := float64(1000000000000000000)

	var bottomRight location
	bottomRightBound := location{x: 100_000, y: 100_000}
	bottomRightSmallestDist := float64(1000000000000000000)

	var candidate float64

	for _, loc1 := range locs {
		candidate = pythag(loc1, topLeftBound)
		if candidate < topLeftSmallestDist {
			topLeftSmallestDist = candidate
			topLeft = loc1
		}

		candidate = pythag(loc1, topRightBound)
		if candidate < topRightSmallestDist {
			topRightSmallestDist = candidate
			topRight = loc1
		}

		candidate = pythag(loc1, bottomLeftBound)
		if candidate < bottomLeftSmallestDist {
			bottomLeftSmallestDist = candidate
			bottomLeft = loc1
		}

		candidate = pythag(loc1, bottomRightBound)
		if candidate < bottomRightSmallestDist {
			bottomRightSmallestDist = candidate
			bottomRight = loc1
		}
	}

	biggest := squareSize(topLeft, bottomRight)
	other := squareSize(topRight, bottomLeft)
	if other > biggest {
		biggest = other
	}

	fmt.Println("part 1:", biggest)
}

func pythag(loc1, loc2 location) float64 {
	return math.Sqrt(math.Pow(float64(loc1.x-loc2.x), 2) + math.Pow(float64(loc1.y-loc2.y), 2))
}

func squareSize(loc1, loc2 location) int {
	width := loc1.x - loc2.x
	if width < 0 {
		width *= -1
	}
	width += 1

	height := loc1.y - loc2.y
	if height < 0 {
		height *= -1
	}
	height += 1

	return width * height
}

func part2() {
}

func main() {
	part1()
	part2()
}
