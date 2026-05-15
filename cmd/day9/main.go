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
	column int
	row    int
}

// Determine these from eyeballing the input
const columnMin = 0
const rowMin = 0
const columnMax = 100_000
const rowMax = 100_000

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
		locs = append(locs, location{column: x, row: y})
	}

	var topLeft location
	topLeftBound := location{column: columnMin, row: rowMin}
	topLeftSmallestDist := math.MaxFloat64

	var topRight location
	topRightBound := location{column: columnMax, row: rowMin}
	topRightSmallestDist := math.MaxFloat64

	var bottomLeft location
	bottomLeftBound := location{column: columnMin, row: rowMax}
	bottomLeftSmallestDist := math.MaxFloat64

	var bottomRight location
	bottomRightBound := location{column: columnMax, row: rowMax}
	bottomRightSmallestDist := math.MaxFloat64

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
	return math.Sqrt(math.Pow(float64(loc1.column-loc2.column), 2) + math.Pow(float64(loc1.row-loc2.row), 2))
}

func squareSize(loc1, loc2 location) int {
	width := loc1.column - loc2.column
	if width < 0 {
		width *= -1
	}
	width += 1

	height := loc1.row - loc2.row
	if height < 0 {
		height *= -1
	}
	height += 1

	return width * height
}

const (
	red   = 1
	green = 2
)

func part2() {
	// Part 2 thinkings
	// first we must build the map.
	// load in the locations of all red tiles
	// record the rows and columns where red tiles appear
	// next we will brute force through all possible squares
	// In order to check if a square would fall inside the boundaries,
	// we just need to check if there is a third red that would make up the
	//

	// ..............
	// ..............
	// ...x####x##x..
	// ..........##..
	// ..........##..
	// ..........xx..
	// ..............
	// ..............
	// ..............

	// doing this for both rows and columns draws our boarder of the problem space.

	// we dont want to "fill in" the boundary inside of tileMap, that would be a lot of redundant space (like potentitally requiring gigabytes of ram to process)

	// once we have a boundary, we need a function that tells us if a given square between two locations is in the boundary or not. If it is in the boundaries and it's the largest so far, record as the largest.
	// repeat until done.

	allLocations := []location{}
	// colRedTiles/rowRedTiles helps us easily access the red tiles for any given column or row
	colRedTiles := make(map[int][]location)
	rowRedTiles := make(map[int][]location)
	tileMap := make(map[location]int)

	for _, line := range strings.Split(input, "\n") {
		xy := strings.Split(line, ",")
		x, err1 := strconv.Atoi(xy[0])
		y, err2 := strconv.Atoi(xy[1])
		if err1 != nil || err2 != nil {
			panic(fmt.Sprintf("bad input data, line: '%s'", line))
		}

		loc := location{column: x, row: y}
		allLocations = append(allLocations, loc)
		colRedTiles[x] = append(colRedTiles[x], loc)
		rowRedTiles[y] = append(rowRedTiles[y], loc)
		tileMap[loc] = red
	}

	for rowNo, redTiles := range rowRedTiles {
		// sanity check, delete me
		for _, tile := range redTiles {
			if tile.row != rowNo {
				panic("logic error earliler, all these should be for the specified row number")
			}
		}
		colMin, colMax := getColumnMinMax(redTiles)
		applyGreenTileColBoundary(colMin, colMax, tileMap)
	}

	for colNo, redTiles := range colRedTiles {
		// sanity check, delete me
		for _, tile := range redTiles {
			if tile.column != colNo {
				panic("logic error earliler, all these should be for the specified col number")
			}
		}
		rowMin, rowMax := getRowMinMax(redTiles)
		applyGreenTileRowBoundary(rowMin, rowMax, tileMap)
	}

	// for any calculated square we can check if it is inside the boundary by:
	// - is one of the other two corners also a red tile?
	// - is the other corner inside the boundary? (are there any green tiles beyond it?)
	// - if both these hold true, we know it is in the boundary

	// printTileMap(tileMap)
	// This looks good. But we should actually be saving the min and max per row for easy access,
	// we dont have to want to go iterating and searching for it. - too hard? too many data points potentially
	currentBiggest := 0
	for _, loc1 := range allLocations {
		for _, loc2 := range allLocations {
			candidateSize := squareSize(loc1, loc2)
			// if this square is potentially the biggest to fit...
			if candidateSize > currentBiggest {
				otherCorner1, otherCorner2 := getOtherCornersOfSquare(loc1, loc2)

				val1, ok1 := tileMap[otherCorner1]
				val2, ok2 := tileMap[otherCorner2]

				if (ok1 && val1 == red && isCornerInTileCluster(otherCorner2, tileMap)) || (ok2 && val2 == red && isCornerInTileCluster(otherCorner1, tileMap)) {
					currentBiggest = candidateSize
				}
			}
		}
	}

	fmt.Println("part 2:", currentBiggest)
}

func isCornerInTileCluster(loc location, tileMap map[location]int) bool {
	// if the corner falls directly on our boundary, then we know its all good!
	_, ok := tileMap[loc]
	if ok {
		return true
	}
	// otherwise its not on a boundary, lets check if going up, down, left, right we hit a tile

	// check up
	upFound := false
	for checkRow := loc.row - 1; checkRow >= rowMin; checkRow-- {
		checkLoc := location{row: checkRow, column: loc.column}
		_, ok := tileMap[checkLoc]
		if ok {
			upFound = true // exit early whenever we can
			break
		}
	}
	if !upFound { // exit early whenever we can
		return false
	}

	// check down
	downFound := false
	for checkRow := loc.row + 1; checkRow <= rowMax; checkRow++ {
		checkLoc := location{row: checkRow, column: loc.column}
		_, ok := tileMap[checkLoc]
		if ok {
			downFound = true // exit early whenever we can
			break
		}
	}
	if !downFound { // exit early whenever we can
		return false
	}

	// check left
	leftFound := false
	for checkCol := loc.column - 1; checkCol >= columnMin; checkCol-- {
		checkLoc := location{row: loc.row, column: checkCol}
		_, ok := tileMap[checkLoc]
		if ok {
			leftFound = true // exit early whenever we can
			break
		}
	}
	if !leftFound { // exit early whenever we can
		return false
	}

	// check right
	rightFound := false
	for checkCol := loc.column + 1; checkCol <= columnMax; checkCol++ {
		checkLoc := location{row: loc.row, column: checkCol}
		_, ok := tileMap[checkLoc]
		if ok {
			rightFound = true // exit early whenever we can
			break
		}
	}
	if !rightFound { // exit early whenever we can
		return false
	}

	return true
}

func getOtherCornersOfSquare(loc1, loc2 location) (location, location) {
	return location{row: loc1.row, column: loc2.column}, location{row: loc2.row, column: loc1.column}
}

// printTileMap is just for use with the sample input to visually confirm i am doing what
// i want to do. The real input is way too big for this to be useful
func printTileMap(tileMap map[location]int) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 14; col++ {
			loc := location{row: row, column: col}
			val, ok := tileMap[loc]
			if !ok {
				fmt.Print(".")
			} else {
				if val == red {
					fmt.Print("#")
				} else if val == green {
					fmt.Print("X")
				}
			}
		}
		fmt.Print("\n")
	}
}

func applyGreenTileColBoundary(minLoc, maxLoc location, tileMap map[location]int) {
	for i := minLoc.column + 1; i < maxLoc.column; i++ {
		loc := location{row: minLoc.row, column: i}
		tileMap[loc] = green
	}
}

func applyGreenTileRowBoundary(minLoc, maxLoc location, tileMap map[location]int) {
	for i := minLoc.row + 1; i < maxLoc.row; i++ {
		loc := location{row: i, column: minLoc.column}
		tileMap[loc] = green
	}
}

func getColumnMinMax(locations []location) (location, location) {
	if len(locations) < 2 {
		panic("you weren't meant to call this with less than 2! EXPLODE")
	}

	var minCol *location
	var maxCol *location

	for _, loc := range locations {
		if minCol == nil || loc.column < minCol.column {
			minCol = &loc
		}
		if maxCol == nil || loc.column > maxCol.column {
			maxCol = &loc
		}
	}

	return *minCol, *maxCol
}

func getRowMinMax(locations []location) (location, location) {
	if len(locations) < 2 {
		panic("you weren't meant to call this with len less than 2! EXPLODE")
	}

	var minRow *location
	var maxRow *location

	for _, loc := range locations {
		if minRow == nil || loc.row < minRow.row {
			minRow = &loc
		}
		if maxRow == nil || loc.row > maxRow.row {
			maxRow = &loc
		}
	}

	return *minRow, *maxRow
}

func main() {
	part1()
	part2()
}
