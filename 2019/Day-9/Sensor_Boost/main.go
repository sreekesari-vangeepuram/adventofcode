package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

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
	fmt.Printf("BOOST keycode with input [%d]: %d\n", 1, runBOOST(Program, 1))
	fmt.Printf("BOOST keycode with input [%d]: %d\n", 2, runBOOST(Program, 2))
}

func runBOOST(Program []int, input int) (keycode int) {
	mem := make([]int, 21200) // 21200 is determined after multiple test-cases
	copy(mem, Program)

	ip, rip := 0, 0         // result [keycode], instruction pointer [ip] and relative pointer [rip]
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
			mem[addr(&mem, ip, C, 1, rip)] = input
			ip += 2

		case 4:
			keycode = mem[addr(&mem, ip, C, 1, rip)]
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
			return keycode

		default:
			log.Fatal("[ERROR]: Unknown opCode encountered!")

		}
	}

	return keycode
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
