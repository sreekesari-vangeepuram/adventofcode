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

func (p *Point) update(direction Point) {
	p.x += direction.x
	p.y += direction.y
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

func printRegistrationIdentifier(grid map[Point]int) {

	var min, max Point
	for position := range grid {
		min = min.Minimum(position)
		max = max.Maximum(position)
	}

	for y := max.y; y >= min.y; y-- {
		for x := min.x; x <= max.x; x++ {
			switch grid[Point{x, y}] {
			case 0:
				fmt.Printf(".")
			case 1:
				fmt.Printf("#")
			default:
				log.Fatal("Unknown color encountered!")
			}
		}
		fmt.Println()
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

	fmt.Printf("Panels painted atleast once: %d\n", len(initPainting(Program, 0)))
	fmt.Println("Registration Identifier:")
	printRegistrationIdentifier(initPainting(Program, 1))

}

func initPainting(Program []int, initialPanel int) map[Point]int {
	directions := map[string]Point{
		"<": Point{-1, 0},
		"^": Point{0, 1},
		">": Point{1, 0},
		"v": Point{0, -1},
	}

	pipeIN := make(chan int, 1)
	pipeOUT := make(chan int)
	SIGTERM := make(chan bool)

	go runProgram(Program, pipeIN, pipeOUT, SIGTERM)

	position, direction := Point{0, 0}, directions["^"]
	panelGrid := make(map[Point]int) // dynamic grid
	panelGrid[position] = initialPanel

	for {

		pipeIN <- panelGrid[position]

		select {
		case color := <-pipeOUT:
			panelGrid[position] = color
			brainSignal := <-pipeOUT
			if brainSignal == 0 /*90* LEFT*/ {
				switch direction {
				case directions["<"]:
					direction = directions["v"]
				case directions["^"]:
					direction = directions["<"]
				case directions[">"]:
					direction = directions["^"]
				case directions["v"]:
					direction = directions[">"]
				}

			} else if brainSignal == 1 /*90* RIGHT*/ {
				switch direction {
				case directions["<"]:
					direction = directions["^"]
				case directions["^"]:
					direction = directions[">"]
				case directions[">"]:
					direction = directions["v"]
				case directions["v"]:
					direction = directions["<"]
				}

			} else {
				log.Fatal("Unkown signal encountered from brain [Intcode]!")
			}

			position.update(direction)

		case <-SIGTERM:
			return panelGrid

		}

	}

}

func runProgram(Program []int, pipeIN <-chan int, pipeOUT chan<- int, SIGTERM chan<- bool) {
	mem := make([]int, 212000) // 21200 * 10 incremented the length of the memory 10 times
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
			SIGTERM <- true

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

	return result
}

func min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}
