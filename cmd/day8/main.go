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

type connection struct {
	loc1     location
	loc2     location
	distance float64
}

type connections struct {
	pending   []connection
	completed []connection
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

		circuitName := uuid.NewSHA1(namespace, []byte(fmt.Sprintf("%d,%d,%d", loc.x, loc.y, loc.z))).String()[:7]

		allCircuits = append(allCircuits, circuitName)

		junctionBoxes[loc] = junctionBox{
			loc:     loc,
			circuit: circuitName,
		}

		circuitToLocation[circuitName] = []location{loc}
	}

	// first, begin by building up our starting list of pendingConnections
	allConnections := connections{
		pending:   []connection{},
		completed: []connection{},
	}

	// for each location, find the closest connection
	for _, junctionBox := range junctionBoxes {
		con := getClosestConnection(junctionBox, junctionBoxes, allConnections)
		allConnections.pending = append(allConnections.pending, con)
	}

	// sort in order of closest to furthest distance
	sort.Slice(allConnections.pending, func(i, j int) bool {
		return allConnections.pending[i].distance < allConnections.pending[j].distance
	})

	connectionCount := 0
	// Note: need to change this from connectionsSample (10) to connectionsInput (1000) when switching between sample.txt and input.txt as puzzle source
	// - sample.txt: "After making the ten shortest connections" so we can check our output with the provided sample...
	// - input.txt: 1000 connections as per instructions
	for connectionCount < connectionsSample {
		// always start by looking at the closest comparison, since they are sorted from lowest to highest distance
		connToAction := allConnections.pending[0]
		allCircuits, circuitToLocation = connectCircuits(allLocations, junctionBoxes, junctionBoxes[connToAction.loc1].circuit, junctionBoxes[connToAction.loc2].circuit, circuitToLocation, allCircuits)
		allConnections = markConnectionActioned(connToAction, allConnections)
		connectionCount += 1

		// now that connToAction has been processed, there are potentially two new closest connections (one from either end of the new connection)
		newConnection1 := getClosestConnection(junctionBoxes[connToAction.loc1], junctionBoxes, allConnections)
		// insert newConnection1 in order
		for i, c := range allConnections.pending {
			if newConnection1.distance < c.distance {
				allConnections.pending = slices.Insert(allConnections.pending, i, newConnection1)
				break
			}
		}
		newConnection2 := getClosestConnection(junctionBoxes[connToAction.loc2], junctionBoxes, allConnections)
		// insert newConnection2 in order
		for i, c := range allConnections.pending {
			if newConnection2.distance < c.distance {
				allConnections.pending = slices.Insert(allConnections.pending, i, newConnection1)
				break
			}
		}
	}

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

func connectCircuits(locations []location, junctionBoxes map[location]junctionBox, circuit1, circuit2 string, cuircuitMap map[string][]location, allCircuits []string) ([]string, map[string][]location) {
	// check if they are already connected and skip in that case
	if circuit1 == circuit2 {
		return allCircuits, cuircuitMap
	}

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

	return allCircuits, cuircuitMap
}

// getClosestConnection will give the closest connection to the startBox, while also filtering out connections to itself
// or to previously found shortest connections that have already been handled
func getClosestConnection(startBox junctionBox, endBoxes map[location]junctionBox, allConnections connections) connection {
	var distance *float64
	var loc1 location
	var loc2 location

	for _, endBox := range endBoxes {
		// skip cases:
		// - finding distance from ourselves
		// - finding distance for two points that already have a pending or completed connection
		// - DONT SKIP SAME CIRCUIT YET!!!
		if (startBox.loc == endBox.loc) || isConnectionAlreadyHandled(startBox.loc, endBox.loc, allConnections) {
			continue
		}

		candidateDistance := straightLineDistance(startBox.loc, endBox.loc)
		if distance == nil || candidateDistance < *distance {
			distance = &candidateDistance
			loc1 = startBox.loc
			loc2 = endBox.loc
		}
	}

	return connection{
		loc1: loc1,
		loc2: loc2,
		// we might need to handle nil pointer in here if we run out of connections to find?
		distance: *distance,
	}
}

func isConnectionAlreadyHandled(loc1, loc2 location, allCoconnections connections) bool {
	for _, conn := range allCoconnections.completed {
		if (conn.loc1 == loc1 && conn.loc2 == loc2) || (conn.loc1 == loc2 && conn.loc2 == loc1) {
			return true
		}
	}
	for _, conn := range allCoconnections.pending {
		if (conn.loc1 == loc1 && conn.loc2 == loc2) || (conn.loc1 == loc2 && conn.loc2 == loc1) {
			return true
		}
	}
	return false
}

func markConnectionActioned(conn connection, allConnections connections) connections {
	for i, c := range allConnections.pending {
		if (c.loc1 == conn.loc1 && c.loc2 == conn.loc2) || (c.loc2 == conn.loc1 && c.loc1 == conn.loc2) {
			actionedConnection := allConnections.pending[i]
			allConnections.completed = append(allConnections.completed, actionedConnection)
			allConnections.pending = slices.Delete(allConnections.pending, i, i+1)
			return allConnections
		}
	}
	panic("logic error: tried to mark a connection actioned that wasn't even pending")
}

func part2() {
}

func main() {
	part1()
	part2()
}
