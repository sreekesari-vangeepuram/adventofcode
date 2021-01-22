package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

/* Globalizing the Program 
   to use all over the code */
var Program []int

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
	for _, s := range bufferArray {
		n, _ := strconv.Atoi(s)
		Program = append(Program, n)
	}

	fmt.Println("Number of affected points by the tractor beam in the next [50x50] point-set:", countAffectedPoint(50, 50))
	fmt.Println(`
(X, Y) = Top-left cornered-point in [100x100] square closet to the emitter.
Then, X * 10,000 + Y =`, findSquare(100, 100))
}

func ok(X, Y int) bool {

	InputSignal   := make(chan int)
	OutputSignal  := make(chan int)
	SIGTERM       := make(chan bool)

	go runProgram(Program, InputSignal, OutputSignal, SIGTERM)

	InputSignal <- X
	InputSignal <- Y
	return 1 == <-OutputSignal
}

func countAffectedPoint(WIDTH, DEPTH int) (acc int) {

	for y := 0; y < DEPTH; y++ {
		for x := 0; x < WIDTH; x++ {
			if !ok(x, y) { continue } // Since more points are outside the beam area!
			acc++
		}
	}; return acc
}

func findSquare(WIDTH, DEPTH int) int {

	var dx, dy, X, Y int = 0, 0, 0, 0
	for/*ever*/ {

		if !ok(X, Y) { X++ }

		dx, dy = X, Y
	ROW: for/*ever*/ {
			if ok(dx, dy) && ok(dx + WIDTH-1, dy) && ok(dx, dy + DEPTH-1) { return dx * 10000 + dy }

			dx++
			if !ok(dx + WIDTH-1, dy) {
				Y++; break ROW
			}
		}

	}
}

func runProgram(Program []int, InputSignal <-chan int, OutputSignal chan<- int, SIGTERM chan<- bool) {
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
			mem[addr(&mem, ip, C, 1, rip)] = <-InputSignal
			ip += 2

		case 4:
			OutputSignal <- mem[addr(&mem, ip, C, 1, rip)]
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
	case 0: // locationition mode
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
