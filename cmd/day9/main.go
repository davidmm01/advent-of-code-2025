package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
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

	// ............oo
	// .xxxxxxxxxxxoo
	// .xxxxxxxxxxxx.
	// .xxxxxxxxxxxx.
	// .xxxxxxxxxxxx.
	// .xxxxxxxxxxxx.
	// .xxxxxxxxxxxx.
	// ooxxxxxxxxxxx.
	// oo............

	// 14, 9

	// ..............
	// ..............
	// ........x.....
	// .x............
	// ..............
	// ..............
	// ....x.........
	// ..............
	// ..............

	// No. simply using our ssquare to see if thats closest to the boundary doesn't work, e.g. look at the above
	// The point at [1, 4] is considered more of a bottom left point than [4,6] by looking at the area of the triangles that get made
	// with it and the corner.
	// Maybe instead we need the straght line distance to the corner?

	var topRight location
	var bottomRight location
	var topLeft location
	var bottomLeft location

	// first iterate over all points to establish the size of the boundaries,
	// and also store all the locations in a slice so that we don't have to recompute these later.
	xMin := 1000000
	xMax := 0
	yMin := 1000000
	yMax := 0

	locs := []location{}

	for _, line := range strings.Split(input, "\n") {
		xy := strings.Split(line, ",")

		x, err1 := strconv.Atoi(xy[0])
		y, err2 := strconv.Atoi(xy[1])
		if err1 != nil || err2 != nil {
			panic(fmt.Sprintf("bad input data, line: '%s'", line))
		}

		if x < xMin {
			xMin = x
		}
		if x > xMax {
			xMax = x
		}
		if y < yMin {
			yMin = y
		}
		if y > yMax {
			yMax = y
		}

		locs = append(locs, location{x: x, y: y})
	}

	topLeftBoundary := location{x: xMin, y: yMin}
	topRightBoundary := location{x: xMax, y: yMin}
	bottomLeftBoundary := location{x: xMin, y: yMax}
	bottomRightBoundary := location{x: xMax, y: yMax}

	topLeftBoundary = location{x: 1000, y: 1000}
	topRightBoundary = location{x: 100_000, y: 1000}
	bottomLeftBoundary = location{x: 1000, y: 100_000}
	bottomRightBoundary = location{x: 100_000, y: 100_000}

	// Next, for each point, see if it is the closest point to the boundary by seeing which point makes the smallest square area between it and the boundary
	topRightSmallest := 1000000000000000000
	bottomRightSmallest := 1000000000000000000
	topLeftSmallest := 1000000000000000000
	bottomLeftSmallest := 1000000000000000000

	for _, loc := range locs {
		candidate := squareSize(topLeftBoundary, loc)
		if candidate < topLeftSmallest {
			topLeftSmallest = candidate
			topLeft = loc
		}

		candidate = squareSize(topRightBoundary, loc)
		if candidate < topRightSmallest {
			topRightSmallest = candidate
			topRight = loc
		}

		candidate = squareSize(bottomLeftBoundary, loc)
		if candidate < bottomLeftSmallest {
			bottomLeftSmallest = candidate
			bottomLeft = loc
		}

		candidate = squareSize(bottomRightBoundary, loc)
		if candidate < bottomRightSmallest {
			bottomRightSmallest = candidate
			bottomRight = loc
		}
	}

	fromTopLeft := squareSize(topLeft, bottomRight)
	fromBottomLeft := squareSize(topRight, bottomLeft)

	fmt.Println("topLeft:", topLeft)
	fmt.Println("bottomRight:", bottomRight)
	fmt.Println("topRight:", topRight)
	fmt.Println("bottomLeft:", bottomLeft)

	var biggest int
	if fromBottomLeft > fromBottomLeft {
		biggest = fromBottomLeft
	} else {
		biggest = fromTopLeft
	}

	fmt.Println("part 1:", biggest)
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

func part1BruteForce() {
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

	var finLoc1 location
	var finLoc2 location

	largest := 0
	for _, loc1 := range locs {
		for _, loc2 := range locs {
			size := squareSize(loc1, loc2)
			if size > largest {
				largest = size
				finLoc1 = loc1
				finLoc2 = loc2
			}
		}
	}
	fmt.Println("part 1 (brute force):", largest)
	fmt.Println("loc1:", finLoc1)
	fmt.Println("loc2:", finLoc2)
}

func part2() {
}

func main() {
	part1BruteForce()
	part1()
	part2()
}
