package main

import (
	_ "embed" // Required for the //go:embed directive
	"fmt"
	"math"
	"slices"
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
	circuit string
}

func straightLineDistance(loc1, loc2 location) float64 {
	return math.Sqrt(
		math.Pow(float64(loc1.x-loc2.x), float64(2)) +
			math.Pow(float64(loc1.y-loc2.y), float64(2)) +
			math.Pow(float64(loc1.z-loc2.z), float64(2)))
}

func part1() {
	// we will make circuit labels against junctionBoxes to be uuids, and deterministicly generated to help
	// the seed will be constructed from the location label
	namespace := uuid.NameSpaceURL

	// read input into a slice of locations
	allLocations := []location{}

	allCircuits := []string{}

	// TODO: we are having loc as both key and a property on the junctionBox, might be handy but idk might delete
	junctionBoxes := make(map[location]junctionBox)

	// help us track the number of locations in each circuit
	circuitToLocation := make(map[string][]location)

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

		allLocations = append(allLocations, loc)

		circuitName := uuid.NewSHA1(namespace, []byte(fmt.Sprintf("%d,%d,%d", loc.x, loc.y, loc.z))).String()

		allCircuits = append(allCircuits, circuitName)

		junctionBoxes[loc] = junctionBox{
			loc:     loc,
			circuit: circuitName,
		}

		circuitToLocation[circuitName] = []location{loc}
	}

	// fmt.Println("initial junction boxes")
	// printJunctionBoxes(allLocations, junctionBoxes)

	// first, begin by building up our starting list of comparisons
	// for each location, find the closest location
	comparisons := []comparison{}
	for _, junctionBox1 := range junctionBoxes {
		cmp := getComparison(junctionBox1, junctionBoxes)
		comparisons = append(comparisons, cmp)
	}

	// sort in order of closest to furthest distance
	sort.Slice(comparisons, func(i, j int) bool {
		return comparisons[i].distance < comparisons[j].distance
	})

	// fmt.Println("initial comparisons")
	// printComparisons(comparisons)

	connections := 0
	// Note: need to change this from connectionsSample (10) to connectionsInput (1000) when switching between sample.txt and input.txt as puzzle source
	// - sample.txt: "After making the ten shortest connections" so we can check our output with the provided sample...
	// - input.txt: 1000 connections as per instructions
	for connections < connectionsSample {
		fmt.Printf("--- new iteration, connections: %d\n", connections)
		// always start by looking at the closest comparison, since they are sorted from lowest to highest distance
		cmp := comparisons[0]
		// fmt.Println("actioning comparison:", cmp)

		// if this a new circuit that must be established
		if junctionBoxes[cmp.loc1].circuit != junctionBoxes[cmp.loc2].circuit {
			connections += 1
			allCircuits = connectCircuits(allLocations, junctionBoxes, junctionBoxes[cmp.loc1].circuit, junctionBoxes[cmp.loc2].circuit, circuitToLocation, allCircuits)
		}

		// if it was already covered by a circuit, that is the same as already being connected, hence we will treat it the same
		// as one that was just connected

		// Determine the next newest closet connection for this JunctionBox, since it might be closer than other pre-calculated results
		newCmp := getComparison(junctionBoxes[cmp.loc1], junctionBoxes)

		// delete comparisons[0] since it has been handled
		comparisons = slices.Delete(comparisons, 0, 1)

		// insert newCmp in order
		for i, cmp := range comparisons {
			if newCmp.distance < cmp.distance {
				comparisons = slices.Insert(comparisons, i, newCmp)
				break
			}
		}

		// fmt.Println("end iteration state:")
		// printJunctionBoxes(allLocations, junctionBoxes)
		// printComparisons(comparisons)
		// printCircuits(allCircuits, circuitToLocation)
		// fmt.Println("\n\n")
	}

	// for circuit, locs := range circuitToLocation {
	// 	fmt.Printf("circuit %s contains: %v\n", circuit, locs)
	// }

	// count up the largest circuits, do this by making a slice of sizes, sort, and math the first 3 elements
	magnitudes := []int{}
	for _, locations := range circuitToLocation {
		magnitude := len(locations)
		inserted := false
		for i := range magnitudes {
			if magnitude > magnitudes[i] {
				magnitudes = slices.Insert(magnitudes, i, magnitude)
				inserted = true
				break
			}
		}
		if !inserted {
			magnitudes = append(magnitudes, magnitude)
		}
	}

	fmt.Println("part 1:", magnitudes[0]*magnitudes[1]*magnitudes[2])
}

func printCircuits(allCircuits []string, cuircuitMap map[string][]location) {
	if len(allCircuits) != len(cuircuitMap) { // sanity check
		fmt.Println("len:", len(allCircuits), len(cuircuitMap))
		fmt.Println("all:", allCircuits)
		fmt.Println("map:", cuircuitMap)
		panic("data error, circuits and mapping out of sync")
	}

	fmt.Println("Circuits")
	for _, circuit := range allCircuits {
		fmt.Printf("%s: %v\n", circuit, cuircuitMap[circuit])
	}
}

func printJunctionBoxes(allLocations []location, junctionBoxes map[location]junctionBox) {
	fmt.Println("JunctionBoxes")
	// note we iter over all locations to print instead of the map, so that we always see them in the same order (helps debugging)
	for _, loc := range allLocations {
		fmt.Printf("loc: (%3d,%3d,%3d) circuit: %s\n", loc.x, loc.y, loc.z, junctionBoxes[loc].circuit)
	}
}

func printComparisons(comparisons []comparison) {
	fmt.Println("Comparisons")
	for _, cmp := range comparisons {
		fmt.Printf("loc1: (%3d,%3d,%3d), loc2: (%3d,%3d,%3d), distance: %5.2f\n", cmp.loc1.x, cmp.loc1.y, cmp.loc1.z, cmp.loc2.x, cmp.loc2.y, cmp.loc2.z, cmp.distance)
	}
}

func connectCircuits(locations []location, junctionBoxes map[location]junctionBox, circuit1, circuit2 string, cuircuitMap map[string][]location, allCircuits []string) []string {
	// fmt.Printf("connecting circuit '%s' to '%s', the latter will be absorbed by the former\n", circuit1, circuit2)

	oldCircuitLocations := cuircuitMap[circuit2]

	// extend circuit 1 with circuit 2 members
	cuircuitMap[circuit1] = append(cuircuitMap[circuit1], oldCircuitLocations...)

	// update the labels for circuit2 members to circuit1
	for _, oldCircuitLoc := range oldCircuitLocations {
		newBox := junctionBoxes[oldCircuitLoc]
		newBox.circuit = circuit1
		junctionBoxes[oldCircuitLoc] = newBox
	}

	// remove circuit2 from map
	delete(cuircuitMap, circuit2)

	// remove circuit2 from allCircuits
	for i := range allCircuits {
		if allCircuits[i] == circuit2 {
			allCircuits = slices.Delete(allCircuits, i, i+1)
			break
		}
	}

	return allCircuits
}

// getComparison will generate a slice of comparison structs. It will give the closest distance for each of the startBoxes to all the endBoxes,
// i.e. can be used to just generate one new comparsion when necessary.
func getComparison(startBox junctionBox, endBoxes map[location]junctionBox) comparison {
	var distance *float64
	var loc1 location
	var loc2 location

	for _, endBox := range endBoxes {
		if (startBox.loc == endBox.loc) || (startBox.circuit == endBox.circuit) {
			// skip cases:
			// - finding distance from ourselves
			// - finding distance for two junction boxes already in the same circuit
			continue
		}

		candidateDistance := straightLineDistance(startBox.loc, endBox.loc)
		if distance == nil || candidateDistance < *distance {
			distance = &candidateDistance
			loc1 = startBox.loc
			loc2 = endBox.loc
		}
	}

	return comparison{
		loc1:     loc1,
		loc2:     loc2,
		distance: *distance,
	}
}

func part2() {
}

func main() {
	part1()
	part2()
}
