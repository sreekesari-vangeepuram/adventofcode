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

var Program []int

func main() {

	if len(os.Args) < 2 {
        log.Fatal(`

[ERROR]: Dataset not supplied!
** Usage: ./main /path/to/file
`)
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

	buff_array := strings.Split(strings.Trim(string(data), "\n "), ",")
	for _, element := range buff_array {
		Program = append(Program, str2int(element))
	}

    part1 := make(chan int)
    part2 := make(chan int)

    go getHighestSignal([]int{3, 1, 2, 4, 0}, part1)
    go getHighestSignal([]int{9, 7, 8, 5, 6}, part2)

    fmt.Println("Highest Signal [ Phase - I ] :", <-part1)
	fmt.Println("Highest Signal [ Phase - II] :", <-part2)
}

func getHighestSignal(phaseSignalElements []int, result chan<- int) {

    permutation := make(chan []int)
    go permutations(phaseSignalElements, permutation)

    bestSignal := 0
    for phaseSettings := range permutation {
        presentSignal := flowAtoE(phaseSettings)
        if presentSignal > bestSignal {
            bestSignal = presentSignal
        }
    }

    result <- bestSignal
}

func flowAtoE(phaseSettings []int) int {
    A := make(chan int, 1)
    B := make(chan int)
    C := make(chan int)
    D := make(chan int)
    E := make(chan int)

    SIGINT := make(chan bool)


    // Concurrently flow the signal through
    // all the amplifiers
    go flowSignal(A, B, SIGINT)
    go flowSignal(B, C, SIGINT)
    go flowSignal(C, D, SIGINT)
    go flowSignal(D, E, SIGINT)
    go flowSignal(E, A, SIGINT)

    // Concurrently send phaseSetting
    // to all the amplifiers
    A <- phaseSettings[0]
    B <- phaseSettings[1]
    C <- phaseSettings[2]
    D <- phaseSettings[3]
    E <- phaseSettings[4]

    // Initializing sequence...
    A <- 0

    // Listening to signal interruption [SIGINT]
    for i := 0; i < 5; i++ {
        <-SIGINT
    }

    return <-A
}

func flowSignal(inputPhaseSignal <-chan int, outputPhaseSignal chan<- int, SIGINT chan<- bool) {

	Intcode := make([]int, len(Program))
	copy(Intcode, Program)

	var opCode, address, p1, p2 int
	ip := 0
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
			Intcode[address] = <-inputPhaseSignal
			ip += 2

		case 4:
			outputPhaseSignal<- Intcode[Intcode[ip+1]]
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
            SIGINT <- true
			break

		}
	}

}

func followRules(Intcode []int, ip int) (opCode, address, p1, p2 int) {

	ins := Intcode[ip] // ins -> instruction
	if ins == 99 {
		return ins, address, p1, p2
	}

	if opCode = ins % 100; numberInSlice([]int{3, 4}, opCode) {

		return opCode, Intcode[ip+1], p1, p2

	} else if numberInSlice([]int{1, 2, 5, 6, 7, 8}, opCode) {

		if !numberInSlice([]int{5, 6}, opCode) {
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

func permutations(arr []int, permChan chan []int) {
    var helper func([]int, int)

    helper = func(arr []int, n int) {

        if n == 1 {

            tmp := make([]int, len(arr))
            copy(tmp, arr)

            permChan <- tmp

        } else {

            for i := 0; i < n; i++{
                helper(arr, n - 1)
                if n % 2 == 1{

                    tmp := arr[i]
                    arr[i] = arr[n - 1]
                    arr[n - 1] = tmp

                } else {

                    tmp := arr[0]
                    arr[0] = arr[n - 1]
                    arr[n - 1] = tmp

                }
            }
        }
    }

    helper(arr, len(arr))
    close(permChan)
}

