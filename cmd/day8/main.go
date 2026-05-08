package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

//go:embed sample.txt
var input string

type comparison struct {
	loc          location
	closest      location
	closestIndex int // the index of the
	distance     float64
}

type location struct {
	x int
	y int
	z int
}

func straightLineDistance(loc1, loc2 location) float64 {
	return math.Sqrt(
		math.Pow(float64(loc1.x-loc2.x), float64(2)) +
			math.Pow(float64(loc1.y-loc2.y), float64(2)) +
			math.Pow(float64(loc1.z-loc2.z), float64(2)))
}

func part1() {
	// read input into a slice of locations
	locationsFromInput := []location{}

	lines := strings.Split(input, "\n")

	for _, line := range lines {
		coords := strings.Split(line, ",")
		errs := []error{}
		x, err := strconv.Atoi(coords[0])
		if err != nil {
			errs = append(errs, err)
		}

		y, err := strconv.Atoi(coords[1])
		if err != nil {
			errs = append(errs, err)
		}

		z, err := strconv.Atoi(coords[2])
		if err != nil {
			errs = append(errs, err)
		}

		if len(errs) != 0 {
			panic(errs)
		}

		locationsFromInput = append(locationsFromInput, location{x: x, y: y, z: z})
	}

	comparisons := []comparison{}

	// for each location, find the closest location
	for _, assessLocation := range locationsFromInput {
		cmp := comparison{
			loc:      assessLocation,
			distance: -1,
		}
		for j, candidateLocation := range locationsFromInput {
			candidateDistance := straightLineDistance(assessLocation, candidateLocation)
			if candidateDistance == 0 {
				continue
			}

			if (candidateDistance != 0 && cmp.distance == -1) || candidateDistance < cmp.distance {
				cmp.distance = candidateDistance
				cmp.closest = candidateLocation
				cmp.closestIndex = j
			}
		}

		comparisons = append(comparisons, cmp)
	}

	// sort in order of closest to furthest distance
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].distance < comparisons[j].distance
	})

	for _, cmp := range comparisons {
		fmt.Println(cmp)
	}

	connections := 0

	circuits := [][]location{}
	fmt.Println("---")
	for _, cmp := range comparisons {
		// Note: need to change this from 10 to 1000 when switching between sample.txt and input.txt as puzzle source
		// - sample.txt: "After making the ten shortest connections" so we can check our output with the provided sample...
		// - input.txt: 1000 connections as per instructions
		if connections == 11 {
			break
		}

		// check if the incoming points are already in a circuit or not (either the first or the second point might already be in a circuit must check both)

		cmpFound := false
		cmpClosestFound := false
		iDest := -1

		for i, circuit := range circuits {
			for _, loc := range circuit {
				if loc == cmp.loc {
					cmpFound = true
					iDest = i
				}

				if loc == cmp.closest {
					cmpClosestFound = true
					iDest = i
				}
			}
		}

		if cmpFound && cmpClosestFound {
			// if they are already in here, don't treat them as a connection, since its the same connection
			continue
		} else {
			connections += 1

			// if neither point was found in a circuit, then add them both to a new circuit
			if !cmpFound && !cmpClosestFound {
				circuits = append(circuits, []location{cmp.loc, cmp.closest})
			}

			if cmpFound && !cmpClosestFound {
				circuits[iDest] = append(circuits[iDest], cmp.closest)
			}

			if !cmpFound && cmpClosestFound {
				circuits[iDest] = append(circuits[iDest], cmp.loc)
			}
		}

		fmt.Println("\ncircuits after processing cmp", cmp)
		for _, circuit := range circuits {
			fmt.Println(circuit)
		}
		fmt.Println("connections:", connections)

	}

	// for debug, to remove
	fmt.Println("end circuits:")
	for _, circuit := range circuits {
		fmt.Println(circuit)
	}
}

func part2() {
}

func main() {
	part1()
	part2()
}
