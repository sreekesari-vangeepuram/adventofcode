package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "regexp"
    "strconv"
)

type Point struct { x, y, z int }

func (p Point) AbsSum() int {
    return abs(p.x) + abs(p.y) + abs(p.z)
}

func (p Point) ScaleUp(p0 Point) Point {
    return Point{
        x: p.x + p0.x,
        y: p.y + p0.y,
        z: p.z + p0.z,
    }
}

func (p Point) ScaleDown(p0 Point) Point {
    return Point{
        x: p.x - p0.x,
        y: p.y - p0.y,
        z: p.z - p0.z,
    }
}

func (p Point) motionDirection() Point {
    return Point{
        x: dir(p.x),
        y: dir(p.y),
        z: dir(p.z),
    }
}

type Moon  struct { pos, vel Point }

func main() {
    if len(os.Args) < 2 {
        log.Fatal(`
[ERROR]: Provide the input dataset!
**Usage: ./main /path/to/file
`)
    }

    StoInt := func(s string) int {
        n, err := strconv.Atoi(s)
        if err != nil { log.Fatal(err) }
        return n
    }

    file, err := os.Open(os.Args[1])
    if err != nil { log.Fatal(err) }

    scanner := bufio.NewScanner(file)
    var lines []string
    for scanner.Scan() { lines = append(lines, scanner.Text()) }
    file.Close()

    pattern := regexp.MustCompile(`^<x=(-?\d+), y=(-?\d+), z=(-?\d+)>$`)
    var moonsData []Moon

    for _, line := range lines {
        axes := pattern.FindStringSubmatch(line)
        x, y, z := StoInt(axes[1]), StoInt(axes[2]), StoInt(axes[3])
        moon := Moon{
            pos: Point{x, y, z},
            vel: Point{0, 0, 0},
        }
        moonsData = append(moonsData, moon)
    }

    c1, c2 := make(chan int, 1), make(chan int, 1)

    go totalEnergy(moonsData, c1)
    go inspectForSameState(moonsData, c2)

    fmt.Println("Total energy of the system after 1000 steps:", <-c1)
    fmt.Printf("History repeated after %d steps!\n", <-c2)
}

func totalEnergy(moonsData []Moon, c1 chan<- int) {

    moons := make([]Moon, len(moonsData))
    copy(moons, moonsData)


    for step := 0; step < 1000; step++ { simulateMoons(moons) }

    var total, pot, kin int = 0, 0, 0
    for _, moon := range moons {
        pot, kin = moon.pos.AbsSum(), moon.vel.AbsSum()
        total += pot * kin
    }

    c1 <- total
}

func inspectForSameState(moonsData []Moon, c2 chan<- int) {

    moons := make([]Moon, len(moonsData))
    copy(moons, moonsData)

    var xs, ys, zs int = 0, 0, 0
    var repeatedState bool = false

    for steps := 1; xs == 0 || ys == 0 || zs == 0; steps++ {
        simulateMoons(moons)

        if xs == 0 {
            repeatedState = true
            for i, moon := range moons {
                if moon.pos.x != moonsData[i].pos.x || moon.vel.x != moonsData[i].vel.x { repeatedState = false; break }
            }

            if repeatedState { xs = steps }
        }

        if ys == 0 {
            repeatedState = true
            for i, moon := range moons {
                if moon.pos.y != moonsData[i].pos.y || moon.vel.y != moonsData[i].vel.y { repeatedState = false; break }
            }

            if repeatedState { ys = steps }
        }

        if zs == 0 {
            repeatedState = true
            for i, moon := range moons {
                if moon.pos.z != moonsData[i].pos.z || moon.vel.z != moonsData[i].vel.z { repeatedState = false; break }
            }

            if repeatedState { zs = steps }
        }

    }

    c2 <- lcm(lcm(xs, ys), zs)
}

func simulateMoons(moons []Moon) {
    for i, m1 := range moons {
        for _, m2 := range moons {
            if m1 == m2 { continue }
            m1.vel = m1.vel.ScaleUp(m2.pos.ScaleDown(m1.pos).motionDirection())
        }
            moons[i] = m1
    }

    for i, moon := range moons { moons[i].pos = moon.pos.ScaleUp(moon.vel) }
}


func abs(n int) int {
    if n < 0 { return -n }
    return n
}

func dir(n int) int {
    if n < 0 { return -1 }
    if n > 0 { return  1 }
    return 0
}

func gcd(n1, n2 int) int {
    for n2 != 0 { n1, n2 = n2, n1 % n2 }
    return n1
}

func lcm(n1, n2 int) int { return n1 / gcd(n1, n2) * n2 }
