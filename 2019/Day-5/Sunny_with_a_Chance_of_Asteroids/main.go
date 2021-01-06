package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
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

	fmt.Println("[DIAGNOSTIC CODES] - I  :", solve(Intcode, 1))
	fmt.Println("[DIAGNOSTIC CODES] - II :", solve(Intcode, 5))
}

func solve(Intcode_copy []int, input int) []int {
	Intcode := make([]int, len(Intcode_copy))
	copy(Intcode, Intcode_copy)

	var opCode, address, p1, p2 int
	outputs := make([]int, 0)

	var ip int = 0
	for {

		opCode, address, p1, p2 = followRules(Intcode, ip)
		switch opCode {
		case 1:
			Intcode[address] = p1 + p2
			ip += 4

		case 2:
			Intcode[address] = p1 * p2
			ip += 4

		case 3:
			Intcode[address] = input
			ip += 2

		case 4:
			outputs = append(outputs, Intcode[Intcode[ip+1]])
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
				Intcode[address] = 1
			} else {
				Intcode[address] = 0
			}
			ip += 4

		case 8:
			if p1 == p2 {
				Intcode[address] = 1
			} else {
				Intcode[address] = 0
			}
			ip += 4

		case 99:
			return outputs

		}
	}

	return outputs
}

func followRules(Intcode []int, ip int) (opCode, address, p1, p2 int) {

	ins := Intcode[ip] // ins -> instruction
	if ins == 99 {
		return ins, address, p1, p2
	}

	if opCode = ins % 100; numberInSlice([]int{3, 4}, opCode) {

		return opCode, Intcode[ip+1], p1, p2

	} else if numberInSlice([]int{1, 2, 5, 6, 7, 8}, opCode) {

		if ! /*not in*/ numberInSlice([]int{5, 6}, opCode) {
			address = Intcode[ip+3]
		}

		if math.Floor(float64(ins%1000)/100) == 1 {
			p1 = Intcode[ip+1]
		} else {
			p1 = Intcode[Intcode[ip+1]]
		}

		if math.Floor(float64(ins%10000)/1000) == 1 {
			p2 = Intcode[ip+2]
		} else {
			p2 = Intcode[Intcode[ip+2]]
		}

		return opCode, address, p1, p2
	}

	return opCode, 0, 0, 0
}

func numberInSlice(slice []int, number int) bool {
	for _, n := range slice {
		if n == number {
			return true
		}
	}

	return false
}
