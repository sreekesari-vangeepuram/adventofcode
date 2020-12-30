package main

import (
    "os"
    "math"
    "io/ioutil"
    "log"
    "strings"
    "strconv"
    "fmt"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("** Usage: ./main <filename>")
    }

    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    str2int := func (r string) int {
        number, err := strconv.Atoi(r)
        if err != nil {
            log.Fatal(err)
        }
        return number
    }

    var masses []int
    buff_array := strings.Split(strings.Trim(string(data), "\n "), "\n")
    for _, str := range buff_array {
        masses = append(masses, str2int(str))
    }

    fmt.Printf("Sum of the fuel requirements: %d\n", part1(masses))
    fmt.Printf("Sum of the fuel requirements for all the modules in the spacecraft: %d\n", part2(masses))

}

func part1(masses []int) int {

    var acc int // accumulator
    for _, number := range masses {
        acc += int(math.Floor(float64(number/3))) - 2
    }

    return acc
}

func part2(masses []int) int {

    var acc int // accumulator
    var buff_num int // buffer for an int
    for _, number := range masses {
        buff_num = int(math.Floor(float64(number/3))) - 2
        for buff_num > 0 {
            acc += buff_num
            buff_num = int(math.Floor(float64(buff_num/3))) - 2
        }
    }
    return acc
}
