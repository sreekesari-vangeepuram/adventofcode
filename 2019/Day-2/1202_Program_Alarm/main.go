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
		log.Fatal("** Usage: ./main <filename>")
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	str2int := func(s string) int {
		number, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		return number
	}

	var Intcode []int
	buff_array := strings.Split(strings.Trim(string(data), "\n "), ",")
	for _, element := range buff_array {
		Intcode = append(Intcode, str2int(element))
	}

	var noun, verb = 12, 2
	fmt.Printf("Value at Intcode[0]: %d\n", part1(Intcode, noun, verb))
	noun, verb, result := part2(Intcode)
	fmt.Printf("100 * %d + %d = %d\n", noun, verb, result)
}

func part1(Intcode_copy []int, noun, verb int) int {
	Intcode := make([]int, len(Intcode_copy))
	copy(Intcode, Intcode_copy)

	Intcode[1] = noun
	Intcode[2] = verb

	for ip := 0; ip <= len(Intcode)-1; ip += 4 /* next_ip = 1 opcode + 3 parameters*/ {

		switch Intcode[ip] {
		case 1:
			Intcode[Intcode[ip+3]] = Intcode[Intcode[ip+1]] + Intcode[Intcode[ip+2]]

		case 2:
			Intcode[Intcode[ip+3]] = Intcode[Intcode[ip+1]] * Intcode[Intcode[ip+2]]

		case 99:
			return Intcode[0]

		}
	}
	return Intcode[0]
}

func part2(Intcode []int) (int, int, int) {
	// output = 19690720
	for noun := 0; noun <= 99; noun++ {
		for verb := 0; verb <= 99; verb++ {
			if part1(Intcode, noun, verb) == 19690720 {
				return noun, verb, 100*noun + verb
			}
		}
	}

	return 0, 0, 0
}
