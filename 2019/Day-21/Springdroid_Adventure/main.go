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

// Springdroid Scripts
var (

/* #Rules: 
 *
 * Springdroid JMP 4 at once
 * Check D whether solid or not
 * Mind the gap (holes between tiles)
 *
 */

WALK = `NOT A J
NOT B T
OR T J
NOT C T
OR T J
AND D J
WALK
`

// Addition of new instructions
RUN = `NOT A J
NOT B T
OR T J
NOT C T
OR T J
AND D J
NOT E T
NOT T T
OR H T
AND T J
RUN
`

)

	message, damage := initSpringDroid(Program, WALK)
	fmt.Println("Message from droid:\n---\n", message, "\n---\n")
	fmt.Println("Amount of hull damage reported :", damage, "\n")

	message, damage  = initSpringDroid(Program, RUN)
	fmt.Println("Message from droid:\n---\n", message, "\n---\n")
	fmt.Println("Amount of hull damage reported :", damage, "\n")

}

func initSpringDroid(Program []int, BotScript string) (string, int) {

	var damage int
	var scriptOutput strings.Builder // Droid Message Buffer

	InputCode  := make(chan int, len(BotScript))
	OutputCode := make(chan int)
	SIGTERM    := make(chan bool)

	go runProgram(Program, InputCode, OutputCode, SIGTERM)

	for _, char := range BotScript { InputCode <- int(char) }

	for {

		select {
		case ASCIIcode := <-OutputCode:
			if ASCIIcode >= 128 {
				damage = ASCIIcode
			} else { scriptOutput.WriteRune(rune(ASCIIcode)) }

		case <-SIGTERM:
			return strings.TrimSpace(scriptOutput.String()), damage
		}
	}
}


func runProgram(Program []int, InputCode <-chan int, OutputCode chan<- int, SIGTERM chan<- bool) {
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
			mem[addr(&mem, ip, C, 1, rip)] = <-InputCode
			ip += 2

		case 4:
			OutputCode <- mem[addr(&mem, ip, C, 1, rip)]
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

