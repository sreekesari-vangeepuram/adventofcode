package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(`

[ERROR]: Please provide the dataset.
** Usage: ./main <path/to/file>

`)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	orbitsData := strings.Split(strings.Trim(string(data), "\n "), "\n")
	orbits := make(map[string]string)
	for _, obj := range orbitsData {
		pair := strings.Split(obj, ")")
		// pair[0] -> parent, pair[1] -> child
		orbits[pair[1]] = pair[0] // since child consists of a single parent
	}

	fmt.Println("Total number of direct and indirect orbits:", part1(orbits))
	fmt.Println("The minimum number of orbital transfers required to move from YOU are orbiting to SAN:", part2(orbits))
}

func part1(orbits map[string]string) (orbitsCount int) {

	var object string
	var ok bool

	for subObject := range orbits {
		for {
			object, ok = orbits[subObject]
			if !ok {
				break
			}
			orbitsCount++
			subObject = object
		}

	}

	return orbitsCount
}

func part2(orbits map[string]string) (orbitalTransferCount int) {
	distanceFrom := make(map[string]int)
	var parent string
	var ok bool

	object, youDistance := orbits["YOU"], 0
	for {
		distanceFrom[object] = youDistance
		parent, ok = orbits[object]

		if !ok {
			break
		}

		youDistance++
		object = parent
	}

	object = orbits["SAN"]
	for {
		youDistance, ok = distanceFrom[object]
		parent = orbits[object]

		if ok {
			orbitalTransferCount += youDistance
			break
		}

		orbitalTransferCount++
		object = parent

	}
	return orbitalTransferCount
}
