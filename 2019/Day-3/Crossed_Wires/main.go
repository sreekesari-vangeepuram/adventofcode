package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "log"
    "strings"
    "strconv"
)

var direction map[string]Point

type Instruction struct {
    direction string
    magnitude int
}

type Point struct {
    x, y int
}

func (p *Point) update(ins Instruction) {
    point := direction[ins.direction]
    if point.x == 0 {
                p.y += point.y * ins.magnitude
    } else {
                p.x += point.x * ins.magnitude
    }
}

func main() {

     direction = map[string]Point {
        "R": Point{+1, 0},
        "L": Point{-1, 0},
        "U": Point{0, +1},
        "D": Point{0, -1},
    }

    if len(os.Args) < 2 {
        log.Fatal("** Usage: ./main <filename>")
    }

    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    buff_array := strings.Split(strings.Trim(string(data), "\n "), "\n")
    wires := make([][]Instruction, len(buff_array))
    for number, wire := range buff_array {
        for _, ins := range strings.Split(wire, ",") {
            num, err := strconv.Atoi(ins[1:])
            if err != nil {
                log.Fatal(err)
            }

            wires[number] = append(wires[number], Instruction{string(ins[0]), num})
        }
        buff_array[number] = ""
    }
    part1, part2 := solve(wires)
    fmt.Printf("Manhattan Distance: %d\n", part1)
    fmt.Printf("Fewest combined steps the wires must take to reach an intersection: %d\n", part2)
}

func intersectionsIn(locus *[][]Point) []Point {
    buff_array := make([]Point, 0)
    for _, point1 := range (*locus)[0] {
        for _, point2 := range (*locus)[1] {
            if (point1 == point2) && (point2 != Point{0, 0}) {
                buff_array = append(buff_array, point1)
            }
        }
    }
    return buff_array
}

func Abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}

func solve(wires [][]Instruction) (int, int) {
    locus := make([][]Point, 2)

    var currentLocation, previousLocation Point
    for number, wire := range wires {
        currentLocation  = Point{0, 0}
        previousLocation = Point{0, 0}
        for i, ins := range wire {
            if i < len(wire)  {
                previousLocation = currentLocation
            }
            currentLocation.update(ins)

            if currentLocation.x - previousLocation.x == 0 {
                if currentLocation.y < previousLocation.y {
                    for y := previousLocation.y; y > currentLocation.y; y-- {
                        locus[number] = append(locus[number], Point{currentLocation.x, y})
                    }
                } else {
                    for y := previousLocation.y; y < currentLocation.y; y++ {
                        locus[number] = append(locus[number], Point{currentLocation.x, y})
                    }
                }
            } else {
                if currentLocation.x < previousLocation.x {
                    for x := previousLocation.x; x > currentLocation.x; x-- {
                        locus[number] = append(locus[number], Point{x, currentLocation.y})
                    }
                } else {
                    for x := previousLocation.x; x < currentLocation.x; x++ {
                        locus[number] = append(locus[number], Point{x, currentLocation.y})
                    }
                }

            }
        }
    }


        var manhattanDistance, min1 int
        intersectionPoints := intersectionsIn(&locus)
        for index, point := range intersectionPoints {
            manhattanDistance = Abs(point.x) + Abs(point.y)
            if index == 0 {
                min1 = manhattanDistance
            }

            if manhattanDistance < min1 {
                min1 = manhattanDistance
            }
        }

        countOfSteps := make([]int, 0)
        var bfn1, bfn2 int
        for _, point := range intersectionPoints {
            bfn1, bfn2 = 0, 0
            for i1, p1 := range locus[0] {
                if p1 == point {
                    bfn1 = i1
                }
            }

            for i2, p2 := range locus[1] {
                if p2 == point {
                    bfn2 = i2
                }
            }

            countOfSteps = append(countOfSteps, bfn1 + bfn2)
        }

        var min2 int = countOfSteps[0]
        for _, stepCount := range countOfSteps {
            if stepCount < min2 {
                min2 = stepCount
            }
        }

    return min1, min2

}
