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
    EMPTY   = iota
    WALL
    BLOCK
    HPADDLE
    BALL
)

type Point struct{ x, y int }

const (
    seqIn = iota
    seqOut
    SIGTERM
)

type Instruction struct { tileType, magnitude int }

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

    fmt.Println("Number of block tiles on the screen when game exits:", countBlockTiles(Program))
    fmt.Println("Player score after the last block is broken:", initGame(Program))
}

func countBlockTiles(Program []int) (blockTilesCount int) {
    input := make(chan int)
    insChan := make(chan Instruction)

    go runProgram(Program, input, insChan)

    tiles := make(map[Point]int)
    var ins Instruction
    var loc Point

    checkType := func(tileType int) {
        if tileType != seqOut { log.Fatal("[ERROR]: Unknown tile type encountered: ", tileType) }
    }

    for {
        ins = <-insChan
        switch ins.tileType {

            case seqOut:
                loc = Point{0, 0}

                loc.x = ins.magnitude
                ins = <-insChan
                checkType(ins.tileType)

                loc.y = ins.magnitude
                ins = <-insChan
                checkType(ins.tileType)

                tiles[loc] = ins.magnitude

            case SIGTERM:
                for _, tile := range tiles { if tile == BLOCK { blockTilesCount++ } }
                return blockTilesCount


            default:
                log.Fatal("[ERROR]: Unknown instruction encountered: ", ins.tileType)
        }
    }

}

func initGame(Program []int) (finalScore int) {

    Program[0] = 2 // 2 Quarters inserted to play the game
    inputChan := make(chan int)
    instructions := make(chan Instruction)

    go runProgram(Program, inputChan, instructions)

    tiles := make(map[Point]int)
    var ins Instruction
    loc, p := Point{0, 0}, Point{-1, 0}

    checkType := func(tileType int) {
        if tileType != seqOut { log.Fatal("[ERROR]: Unknown tile type encountered: ", tileType) }
    }

    for {

        ins = <-instructions
        switch ins.tileType {

            case seqIn:
                var hPaddle, ball Point
                for position, tile := range tiles {
                    switch tile {
                        case HPADDLE:
                            hPaddle = position
                        case BALL:
                            ball = position
                    }
                }

                inputChan <- joystickPosition(ball.x - hPaddle.x)

            case seqOut:
                loc = Point{0, 0}

                loc.x = ins.magnitude
                ins = <-instructions
                checkType(ins.tileType)

                loc.y = ins.magnitude
                ins = <-instructions
                checkType(ins.tileType)

                if loc == p {
                    finalScore = ins.magnitude
                } else {
                    tiles[loc] = ins.magnitude
                }

            case SIGTERM:
                return finalScore

            default:
                log.Fatal("[ERROR]: Unknown instruction encountered: ", ins.tileType)
        }
    }
}

func runProgram(Program []int, input <-chan int, instructions chan<- Instruction) {
	mem := make([]int, len(Program)) // Expandable capacity, checkout `addr` func
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
            instructions <- Instruction{ tileType: seqIn }
			mem[addr(&mem, ip, C, 1, rip)] = <-input
			ip += 2

		case 4:
			instructions <- Instruction{tileType: seqOut, magnitude: mem[addr(&mem, ip, C, 1, rip)]}
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
			instructions <- Instruction{ tileType: SIGTERM }
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

func joystickPosition(n int) int {
    if n < 0 { return -1 }
    if n > 0 { return  1 }
    return 0
}

