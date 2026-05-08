package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

//go:embed sample.txt
var input string

const connectionsSample = 10
const connectionsInput = 1000

type comparison struct {
	loc1     location
	loc2     location
	distance float64
}

type location struct {
	x int
	y int
	z int
}

type junctionBox struct {
	loc     location
	circuit string // "" == individuial circuit
	// nextClosest *location // pointer to differentiate between unset and set
	// distance    *float64  // pointer to differentiate between unset and set, else default value of 0 will require special logic around it
}

func straightLineDistance(loc1, loc2 location) float64 {
	return math.Sqrt(
		math.Pow(float64(loc1.x-loc2.x), float64(2)) +
			math.Pow(float64(loc1.y-loc2.y), float64(2)) +
			math.Pow(float64(loc1.z-loc2.z), float64(2)))
}

// algorithm:
// allLocations := list of the locations of juction boxes. this is necessary to iterate over all juction boxes, since they are being stored as a map (not true)
// with
//

func part1() {
	// we will make circuit labels against junctionBoxes to be uuids, and deterministicly generated to help
	// the seed will be constructed from the location label
	namespace := uuid.NameSpaceURL

	// read input into a slice of locations
	allLocations := []location{}
	// TODO: we are having loc as both key and a property on the junctionBox, might be handy but idk might delete
	junctionBoxes := make(map[location]junctionBox)

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

		loc := location{x: x, y: y, z: z}
		junctionBoxes[loc] = junctionBox{
			loc:     loc,
			circuit: uuid.NewSHA1(namespace, []byte(fmt.Sprintf("%d,%d,%d", loc.x, loc.y, loc.z))).String(),
		}

		allLocations = append(allLocations, loc)
	}

	// first, begin by building up our starting list of comparisons
	// for each location, find the closest location
	comparisons := []comparison{}

	for _, junctionBox1 := range junctionBoxes {
		var distance *float64
		var loc1 location
		var loc2 location

		for _, junctionBox2 := range junctionBoxes {
			if (junctionBox1.loc == junctionBox2.loc) || ((junctionBox1.circuit == junctionBox2.circuit) && (junctionBox1.circuit != "" && junctionBox2.circuit != "")) {
				// skip cases:
				// - finding distance from ourselves
				// - finding distance for two junction boxes already in the same circuit
				continue
			}

			candidateDistance := straightLineDistance(junctionBox1.loc, junctionBox2.loc)
			if distance == nil || candidateDistance < *distance {
				distance = &candidateDistance
				loc1 = junctionBox1.loc
				loc2 = junctionBox2.loc
			}
		}

		comparisons = append(comparisons, comparison{
			loc1:     loc1,
			loc2:     loc2,
			distance: *distance,
		})
	}

	// sort in order of closest to furthest distance
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].distance < comparisons[j].distance
	})

	connections := 0
	// Note: need to change this from connectionsSample (10) to connectionsInput (1000) when switching between sample.txt and input.txt as puzzle source
	// - sample.txt: "After making the ten shortest connections" so we can check our output with the provided sample...
	// - input.txt: 1000 connections as per instructions
	for connections < connectionsSample {
		fmt.Println("---")
		// always start by looking at the closest comparison, since they are sorted from lowest to highest distance
		cmp := comparisons[0]

		// is this comparison already covered by a circuit?
		if junctionBoxes[cmp.loc1].circuit == junctionBoxes[cmp.loc2].circuit {
			todo
		} else { // not already covered by a circuit, make a new connection
			todo
		}

		// Determine the next newest closet connection for this JunctionBox

	}
}

func part2() {
}

func main() {
	part1()
	part2()
}
