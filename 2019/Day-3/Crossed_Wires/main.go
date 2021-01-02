package main

import (
    "os"
    "log"
    "io/ioutil"
    "fmt"
    "strings"
    "strconv"
)

type Instruction struct {
    dir string
    mag int
}

func parseInstructions(instructions []string) []Instruction {
    buff_array := make([]Instruction, len(instructions))
    for i, ins := range instructions {
        mag, _ := strconv.Atoi(string(ins[1:]))
        buff_array[i] = Instruction{string(ins[0]), mag}
    }
    return buff_array
}

type Point struct {
    x, y int
}

var direction map[string]Point
func (p *Point) update(ins Instruction) {
    point := direction[ins.dir]
    if point.x == 0 {
                p.y += point.y * ins.mag
    } else {
                p.x += point.x * ins.mag
    }
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("** Usage: ./main <filename>")
    }

    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    var buff_array [][]string
    for _, line := range strings.Split(strings.Trim(string(data), "\n "), "\n") {
        buff_array = append(buff_array, strings.Split(line, ","))
    }

    // Parse Instructions
    wire1 := parseInstructions(buff_array[0])
    wire2 := parseInstructions(buff_array[1])

    c1 := make(chan Point, 1)
    c2 := make(chan Point, 1)

    direction = map[string]Point {
        "R": Point{+1, 0},
        "L": Point{-1, 0},
        "U": Point{0, +1},
        "D": Point{0, -1},
    }

    // Fetching tracks with the 
    // obtained instructions `wire1`, `wire2`
    var wire1Track, wire2Track []Point

    go fetchTrack(wire1, c1)
    for point1 := range c1 {
        wire1Track = append(wire1Track, point1)
    }

    go fetchTrack(wire2, c2)
    for point2 := range c2 {
        wire2Track = append(wire2Track, point2)
    }

    var intersectionPoints []Point
    /* Ignoring the 1st intersecting 
     *  point at Point{0, 0}
     */
    for _, p1 := range wire1Track[1:] {
        for _, p2 := range wire2Track[1:] {
            if p1 == p2 {
                intersectionPoints = append(intersectionPoints, p1)
            }
        }
    }

    var manhattanDistance, bfn1, bfn2 int
    var countOfSteps []int

    // PART - 1
    part1 := int(^uint(0) >> 1)
    for _, point := range intersectionPoints {
        manhattanDistance = Abs(point.x) + Abs(point.y)
        if manhattanDistance < part1 {
            part1 = manhattanDistance
        }
    }

    // PART - 2
    for _, point := range intersectionPoints {
        bfn1, bfn2 = 0, 0
        for i1, p1 := range wire1Track {
            if p1 == point {
                bfn1 = i1
            }
        }

        for i2, p2 := range wire2Track {
            if p2 == point {
                bfn2 = i2
            }
        }

        countOfSteps = append(countOfSteps, bfn1 + bfn2)
    }

    part2 := countOfSteps[0]
    for _, stepCount := range countOfSteps {
        if stepCount < part2 {
            part2 = stepCount
        }
    }

    fmt.Printf("Manhattan Distance: %d\n", part1)
    fmt.Printf("Fewest combined steps the wires must take to reach an intersection: %d\n", part2)
}

func fetchTrack(instructions []Instruction, locus chan Point) {
    currentLocation  := Point{0, 0}
    previousLocation := Point{0, 0}
    for i, instruction := range instructions {
        if i < len(instructions)  {
            previousLocation = currentLocation
        }

        currentLocation.update(instruction)

        if currentLocation.x - previousLocation.x == 0 {
            if currentLocation.y < previousLocation.y {
                for y := previousLocation.y; y > currentLocation.y; y-- {
                    locus <- Point{currentLocation.x, y}
                }
            } else {
                for y := previousLocation.y; y < currentLocation.y; y++ {
                    locus <- Point{currentLocation.x, y}
                }
            }
        } else {
            if currentLocation.x < previousLocation.x {
                for x := previousLocation.x; x > currentLocation.x; x-- {
                    locus <- Point{x, currentLocation.y}
                }
            } else {
                for x := previousLocation.x; x < currentLocation.x; x++ {
                    locus <- Point{x, currentLocation.y}
                }
            }

        }
    }
    close(locus)

}

func Abs(number int) int {
    if number < 0 {
        return -number
    }
    return number
}

