package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
    WALL    = iota
    PATH
    FOUND
)

const (
    NORTH   = iota + 1
    SOUTH
    WEST
    EAST
)

type Point struct { x, y int }

func (p Point) Scale(p0 Point) Point {
    return Point{
        x: p.x + p0.x,
        y: p.y + p0.y,
    }
}

func (p Point) Minimum(p0 Point) Point {
    return Point{
        x: min(p.x, p0.x),
        y: min(p.y, p0.y),
    }
}

func (p Point) Maximum(p0 Point) Point {
    return Point{
        x: max(p.x, p0.x),
        y: max(p.y, p0.y),
    }
}

// Queue data structure
type Track struct {
    position    Point
    distance    int
    next        *Track
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal(`

[ERROR]: Provide the dataset!
**Usage: ./main /path/to/file
`)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	bufferArray := strings.Split(strings.Trim(string(data), " \n"), ",")
	Program := make([]int, len(bufferArray))

	for i, s := range bufferArray {
		n, _ := strconv.Atoi(s)
		Program[i] = n
	}

	oxygenSystemLocation, minimumReach, time := initRepairDroid(Program)
	fmt.Printf("Fewest number of commands to reach Oxygen System at %v: %d\n", oxygenSystemLocation, minimumReach)
	fmt.Printf("Time took to fill the oxygen tank-2 with oxygen: %d minutes\n", time)
}

// Some initializers
var (
    UP      = Point{ 0, -1}
    DOWN    = Point{ 0,  1}
    LEFT    = Point{-1,  0}
    RIGHT   = Point{ 1,  0}

    directions  = [4]Point{UP, DOWN, LEFT, RIGHT}
    inputs      = [4]int{NORTH, SOUTH, WEST, EAST}
    reverse     = map[int]int{NORTH: SOUTH, SOUTH: NORTH, WEST: EAST, EAST: WEST}
    direction   = map[int]Point{NORTH: UP, SOUTH: DOWN, WEST: LEFT, EAST: RIGHT}
)

func explorePath(grid map[Point]int, location, target Point, inputSignal, outputSignal chan int) Point {
	var cell *Track

	var track []Track
	track = append(track, Track{position: target, distance: 0})
	visited := make(map[Point]bool)
	visited[target] = true

	for len(track) != 0 {

		node := track[0]
		track = track[1:]

		if node.position == location { cell = &node; break }

		for _, dir := range direction {
			next := node.position.Scale(dir)
			if grid[next] == PATH && !visited[next] {
				track = append(track, Track{position: next, distance: node.distance + 1, next: &node})
				visited[next] = true
			}
		}

	}

	for cell.next != nil {
		for command, dir := range direction {
			if cell.position.Scale(dir) == cell.next.position {
				inputSignal <- command
				<-outputSignal
				break
			}
		}

		cell = cell.next
	}

	return target
}

func initRepairDroid(Program []int) (Point, int, int) {

	inputSignal, outputSignal, SIGTERM := make(chan int), make(chan int), make(chan bool)

	go runProgram(Program, inputSignal, outputSignal, SIGTERM)

	loc  := Point{0, 0}
	grid := make(map[Point]int)
	grid[loc] = PATH

	var OxySysDistance int
	var OxySysLocation Point

	var track []Track
	track = append(track, Track{position: loc, distance: 0})

	for len(track) != 0 {

			node := track[0]
			track = track[1:]
			loc = explorePath(grid, loc, node.position, inputSignal, outputSignal)

			for _, command := range inputs {
				next, nextDist := loc.Scale(direction[command]), node.distance + 1
				_, ok := grid[next]
				if !ok {

					inputSignal <- command
					switch <-outputSignal {
					case WALL:
							grid[next] = WALL

					case FOUND:
							if OxySysDistance == 0 {OxySysDistance, OxySysLocation = nextDist, next}
							fallthrough

					case PATH:
							grid[next] = PATH
							track = append(track, Track{position: next, distance: nextDist})

							inputSignal <- reverse[command]
							<-outputSignal
					}
				}
			}
	}

	// Intiatlize Timer
	var time int
	var navTrack []Track
	navTrack = append(navTrack, Track{position: OxySysLocation, distance: 0})

	visited := make(map[Point]bool)
	visited[OxySysLocation] = true	// navigating from oxygen system location

	for len(navTrack) != 0 {
			node := navTrack[0]
			navTrack = navTrack[1:]

			time = max(node.distance, time)		// Since time is calculated as O2 volume expands

			for _, dir := range directions {
				next := node.position.Scale(dir)
				if grid[next] == PATH &&  !visited[next] {
					navTrack = append(navTrack, Track{position: next, distance: node.distance + 1})
					visited[next] = true
				}
			}
	}

	return OxySysLocation, OxySysDistance, time
}

func runProgram(Program []int, inputSignal <-chan int, outputSignal chan<- int, SIGTERM chan<- bool) {
	mem := make([]int, len(Program)) // Expandable in `addr` func, on need only!
	copy(mem, Program)

	ip, rip := 0, 0         // instruction pointer [ip] and relative pointer [rip]
	var opCodeObject [5]int // array to split ABCDE
	var A, B, C, DE int     // digits: A, B, C in `ABCDE`
	var buffNum, p1, p2 int // buffer integer variables

	for {

		buffNum = int(mem[ip])
		for i := 4; i >= 0; i-- {
			opCodeObject[i] = buffNum % 10
			buffNum /= 10
		}

		DE = opCodeObject[3]*10 + opCodeObject[4] // opCode `DE` in `ABCDE` is extracted
		A, B, C = opCodeObject[0], opCodeObject[1], opCodeObject[2]

		p1, p2 = mem[addr(&mem, ip, C, 1, rip)], mem[addr(&mem, ip, B, 2, rip)]
		switch DE {

		case 1:
			mem[addr(&mem, ip, A, 3, rip)] = p1 + p2
			ip += 4

		case 2:
			mem[addr(&mem, ip, A, 3, rip)] = p1 * p2
			ip += 4

		case 3:
			mem[addr(&mem, ip, C, 1, rip)] = <-inputSignal
			ip += 2

		case 4:
			outputSignal <- mem[addr(&mem, ip, C, 1, rip)]
			ip += 2

		case 5:
			if p1 != 0 {
				ip = p2
			} else {
				ip += 3
			}

		case 6:
			if p1 == 0 {
				ip = p2
			} else {
				ip += 3
			}

		case 7:
			if p1 < p2 {
				mem[addr(&mem, ip, A, 3, rip)] = 1
			} else {
				mem[addr(&mem, ip, A, 3, rip)] = 0
			}
			ip += 4

		case 8:
			if p1 == p2 {
				mem[addr(&mem, ip, A, 3, rip)] = 1
			} else {
				mem[addr(&mem, ip, A, 3, rip)] = 0
			}
			ip += 4

		case 9:
			rip += mem[addr(&mem, ip, C, 1, rip)]
			ip += 2

		case 99:
			SIGTERM <- true
            return

		default:
			log.Fatal("[ERROR]: Unknown opCode encountered!")

		}
	}
}

func addr(mem *[]int, ip int, digit int, offset int, rip int) (result int) {

	switch digit {
	case 0: // position mode
		result = (*mem)[ip+offset]

	case 1: // immediate mode
		result = ip + offset

	case 2: // relative mode
		result = (*mem)[ip+offset] + rip
	}

    // Expands the capacity on need
    for len(*mem) <= result { *mem = append(*mem, 0) }

	return abs(result)
}

func abs(n int) int {
    if n < 0 { return -n }
    return n
}

func min(n1, n2 int) int {
    if n1 < n2 { return n1 }
    return n2
}

func max(n1, n2 int) int {
    if n1 > n2 { return n1 }
    return n2
}
