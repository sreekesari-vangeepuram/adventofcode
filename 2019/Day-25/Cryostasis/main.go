package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)


func main() {

	if len(os.Args) < 2 {
		log.Fatal(`
[ERROR]: Provide the dataset!
**Usage: ./main /path/to/file
`)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil { log.Fatal(err) }

	bufferArray := strings.Split(strings.Trim(string(data), " \n"), ",")
	Program := make([]int, len(bufferArray))
	for i, s := range bufferArray {
		n, _ := strconv.Atoi(s)
		Program[i] = n
	}

	fmt.Println(discoverPassword(Program))
}

func discoverPassword(Program []int) (int, string) {

	STDIN   := make(chan int)
	STDOUT  := make(chan int)
	SIGTERM := make(chan bool)

	go runProgram(Program, STDIN, STDOUT, SIGTERM)

	// User Input reader
	go func() {
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
				for _, c := range reader.Bytes() { STDIN <- int(c) }
				STDIN <- int('\n') // ASCII 10
		}
	}()

	var bufferBuilder strings.Builder
	for/*ever*/ {
		select {
			case char := <-STDOUT:
				switch {
				case char == '\n':
					fmt.Println(bufferBuilder.String())
					bufferBuilder.Reset()

				case char >= 0 && char <= unicode.MaxASCII:
					bufferBuilder.WriteByte(byte(char))

				default:
					fmt.Println(bufferBuilder.String())
					return char, "DONE"
				}

			case <-SIGTERM:
				return 0, "INTCODE HALT NOW!" // Process terminated...
		}
	}

	return 1, "INTCODE CRASHED!"
}

func runProgram(Program []int, STDIN <-chan int, STDOUT chan<- int, SIGTERM chan<- bool) {
	mem := make([]int, len(Program)) // Expandable in `addr` func, on need only!
	copy(mem, Program)

	var ip, rip int = 0, 0  // instruction pointer [ip] and relative pointer [rip]
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
			mem[addr(&mem, ip, C, 1, rip)] = <-STDIN
			ip += 2

		case 4:
			STDOUT <- mem[addr(&mem, ip, C, 1, rip)]
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
