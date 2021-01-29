package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Point struct{ x, y int }

func (p Point) ScaleUp(p0 Point) Point {
	return Point{
		x:	p.x + p0.x,
		y:	p.y + p0.y,
	}
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

	view, acc := captureTheView(Program), 0

	for x := 1; x+1 < len(view); x++ {
		for y := 1; y+1 < len(view[0]); y++ {
			if view[x][y] == '#' {
				if view[x+1][y] == '#' && view[x-1][y] == '#' {
					if view[x][y+1] == '#' && view[x][y-1] == '#' {
						acc += x * y
					}
				}
			}
		}
	}


	fmt.Println("Sum of the alignment parameters for the scaffold intersections:", acc)
	fmt.Println("Number of dust particles collected by vaccum robot:", collectAndCountDust(Program, view))
}

func captureTheView(Program []int) []string {

	pipeIN := make(chan int)
	pipeOUT := make(chan int)
	SIGTERM := make(chan bool)

	go runProgram(Program, pipeIN, pipeOUT, SIGTERM)

	var accumulator strings.Builder

	LOOP: for {

		select {
		case ASCIIcode := <-pipeOUT:
			accumulator.WriteRune(rune(ASCIIcode))

		case <-SIGTERM:
			break LOOP	// breaking the `LOOP` with label
		}
	}

	// Assemble all the units
	return strings.Split(strings.TrimSpace(accumulator.String()), "\n")
}

var (

	UP    = Point{0, -1}
	DOWN  = Point{0, 1}
	LEFT  = Point{-1, 0}
	RIGHT = Point{1, 0}

	TL = map[Point]Point{UP: LEFT, LEFT: DOWN, DOWN: RIGHT, RIGHT: UP}
	TR = map[Point]Point{UP: RIGHT, RIGHT: DOWN, DOWN: LEFT, LEFT: UP}
)

func shortenedRoute(track []string, fns [][]string, frags [][]string) (result [][4][]string) {

	if len(fns) == 3 {

		if len(frags) != 0 { return nil }

		// Compute mainfn
		var mainfn []string
		for len(track) != 0 {
			for i, fn := range fns {
				if HasPrefix(track, fn) {
					track = track[len(fn):]
					mainfn = append(mainfn, string('A'+i))
				}
			}
		}

		// Check mem overflow for mainfn
		if len(strings.Join(mainfn, ",")) > 20 { return nil	}

		var Prog [4][]string
		Prog[0] = mainfn
		Prog[1] = fns[0]
		Prog[2] = fns[1]
		Prog[3] = fns[2]

		result = append(result, Prog)
		return result
	}

	// Use first fragment
	fragment := frags[0]
	var blocks [][]string

	for trackLen := 1; trackLen <= len(fragment); trackLen++ {
		block := fragment[:trackLen]
		if len(strings.Join(block, ",")) <= 20 {
			blocks = append(blocks, block)
		}
	}

	// Done! Send the result...
	if len(frags) == 0 {
		newfns := make([][]string, 0, 4)
		newfns = append(newfns, fns...)
		newfns = append(newfns, []string{})

		subResult := shortenedRoute(track, newfns, frags)
		result = append(result, subResult...)
		return result
	}


	// Verify another block
	for _, block := range blocks {
		var newfgs [][]string
		for _, fragment := range frags {
			for {
				i := Index(fragment, block)
				if i == -1 { break }
				if i != 0 {
					newfgs = append(newfgs, fragment[:i])
				};	fragment = fragment[len(block) + i:]
			}

			if len(fragment) != 0 {	newfgs = append(newfgs, fragment) }
		}

		newfns := make([][]string, 0, 4)
		newfns = append(newfns, fns...)
		newfns = append(newfns, block) // Block -> fns.

		subResult := shortenedRoute(track, newfns, newfgs)
		result = append(result, subResult...)
	}

	return result
}

func collectAndCountDust(Program []int, view []string) int {

		var track []string

		Program[0] = 2  // Init vaccum robot...

		// Locate the robot!
		var loc, dir Point
		for y := 0; y < len(view); y++ {

			for x := 0; x < len(view[0]); x++ {

				switch view[y][x] {
				case '<':
					dir, loc = LEFT, Point{x, y}

				case '^':
					dir, loc = UP, Point{x, y}

				case '>':
					dir, loc = RIGHT, Point{x, y}

				case 'v':
					dir, loc = DOWN, Point{x, y}

				}
			}
		}

		theBlockIsScaffold := func(loc Point) bool {
				cX := loc.x >= 0 && loc.x < len(view[0])
				cY := loc.y >= 0 && loc.y < len(view)
			return cX && cY && view[loc.y][loc.x] == '#'
		}

		for {

			trackLen := 0
			for theBlockIsScaffold(loc.ScaleUp(dir)) {
				loc = loc.ScaleUp(dir); trackLen++
			}

			if trackLen != 0 { track = append(track, strconv.Itoa(trackLen)) }

			turn90Left, turn90Right := TL[dir], TR[dir]
			if theBlockIsScaffold(loc.ScaleUp(turn90Left)) {
				dir, track = turn90Left, append(track, "L")
			} else if theBlockIsScaffold(loc.ScaleUp(turn90Right)) {
				dir, track = turn90Right, append(track, "R")
			} else { break }
		}

		// Buffer size 72 is enough, but 4 * 20 is used instead
		pipeIN := make(chan int, 80)
		pipeOUT := make(chan int)
		SIGTERM := make(chan bool)

		go runProgram(Program, pipeIN, pipeOUT, SIGTERM)

		totalRoute :=shortenedRoute(track, nil, [][]string{track})
		if len(totalRoute) == 0 { log.Fatal("[ERROR]: Totally filtered!") }

		fns := totalRoute[0]

		var mainfn, fnA, fnB, fnC string // Path routines
		mainfn	= strings.Join(fns[0], ",")
		fnA		= strings.Join(fns[1], ",")
		fnB		= strings.Join(fns[2], ",")
		fnC		= strings.Join(fns[3], ",")

		route := fmt.Sprintf("%s\n%s\n%s\n%s\nn\n", mainfn, fnA, fnB, fnC)
		for _, c := range route { pipeIN <- int(c) }

	FINAL_LOOP:	for {
			select {
			case ASCIIcode := <-pipeOUT:
				if ASCIIcode >= 128 {
					return ASCIIcode
				}

			case <-SIGTERM:
				break FINAL_LOOP
			}
		}

	return 0

}

func runProgram(Program []int, pipeIN <-chan int, pipeOUT chan<- int, SIGTERM chan<- bool) {
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
			mem[addr(&mem, ip, C, 1, rip)] = <-pipeIN
			ip += 2

		case 4:
			pipeOUT <- mem[addr(&mem, ip, C, 1, rip)]
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
			SIGTERM <- true; return

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

func HasPrefix(slice, prefix []string) (result bool) {
	if  len(prefix) > len(slice) { return false }
	for i, s := range prefix {
		if slice[i] != s { return false }
	}; return true
}

func Index(slice, subSlice []string) int {
	for i := 0; i <= len(slice)-len(subSlice); i++ {
		if HasPrefix(slice[i:], subSlice) {	return i}
	};  return -1
}
