package main

import (
    "os"
    "io/ioutil"
    "log"
    "strings"
    "fmt"
    "math"
    "sort"
)

type Point struct { x, y int }

// Method for distance between 2 points
func (p Point) dist(p0 Point) float64 {
    return math.Sqrt(math.Pow(float64(p.x - p0.x), 2) + math.Pow(float64(p.y - p0.y), 2))
}

func main() {

    if len(os.Args) < 2 {
        log.Fatal(`

[ERROR]: Provide the dataset!
**Usage: ./main /path/to/file
`)
    }

    data, err := ioutil.ReadFile(os.Args[1])
    if err != nil { log.Fatal(err) }

    grid := strings.Split(strings.Trim(string(data), " \n"), "\n")
    asteroidsMap := make(map[Point]bool)
    asteroids := make([][]int, len(grid))

    for y, row := range grid {
        asteroids[y] = make([]int, len(grid[y]))
        for x, block := range row {
            if block == '#' {
                asteroidsMap[Point{x, y}] = true
                asteroids[y][x] = 1
            } else if block != '.' {
                log.Fatal("Unknown symbol[%c] encountered!", block)
            }
        }
    }

    bestCount, bestLocation := findBestLocation(asteroidsMap)
    fmt.Printf("Number of asteroids encountered from best location %v : %d\n", bestLocation, bestCount)

    asteroid200 := vaporize(bestLocation, asteroids, 200)

    fmt.Printf(`Location of 200th asteroid to vaporize: %v
asteroid_200[X] * 100 + asteroid_200[Y] = %d
`, asteroid200, asteroid200.x * 100 + asteroid200.y)

}


func findBestLocation(asteriods map[Point]bool) (int, Point) {

    var bestCount int
    var dx, dy int
    var bestLocation Point

    for p1 := range asteriods {
        Set := make(map[Point]bool)
        for p2 := range asteriods {
            if p1 != p2 {
                    dx, dy = p2.x - p1.x, p2.y - p1.y
                    dx, dy = dx/HCF(dx, dy), dy/HCF(dx, dy)
                    Set[Point{dx, dy}] = true
                }
            }

            if len(Set) > bestCount {
                bestCount = len(Set)
                bestLocation = p1
        }
    }
    return bestCount, bestLocation
}



func vaporize(station Point, asteroids [][]int, targetAsteroidNumber int) Point {

    // Count and map the slopes to asteroids
    dirs := make(map[Point][]Point) // Slopes
    for y := 0; y < len(asteroids); y++ {
        for x := 0; x < len(asteroids[0]); x++ {
            currentPoint := Point{x: x, y: y}
            if asteroids[y][x] == 1 && currentPoint != station {

                dir := Point{
                    x: x - station.x,
                    y: y - station.y,
                }
                // Scale down
                dir = Point{
                    x: dir.x/HCF(dir.x, dir.y),
                    y: dir.y/HCF(dir.x, dir.y),
                }

                dirs[dir] = append(dirs[dir], currentPoint)
            }
        }
    }

    getAngleOf := func (slope Point) (angle float64) {
        angle = math.Atan2(float64(slope.y), float64(slope.x))/(2.0 * math.Pi) * 360.0 + 90.0 
        if angle < 0 { angle += 360.0 }
        return angle
    }

    angles := make([]float64, 0)
    asteroidAtAngle := make(map[float64][]Point)

    for slope := range dirs {

        angle := getAngleOf(slope)
        asteroidAtAngle[angle] = dirs[slope]

        sort.Slice(asteroidAtAngle[angle], func(i, j int) bool {
            return asteroidAtAngle[angle][i].dist(station) < asteroidAtAngle[angle][j].dist(station)
        })

        if len(asteroidAtAngle[angle]) >= 1 { angles = append(angles, angle) }
        sort.Float64s(angles)

    }

    for i := 0; i < targetAsteroidNumber; {
        for _, angle := range angles {

            if len(asteroidAtAngle[angle]) > 0 {
                if i ++; targetAsteroidNumber == i {
                    return asteroidAtAngle[angle][0]// Asteroid location will be found here!
                }

                asteroidAtAngle[angle] = asteroidAtAngle[angle][1:] // filter
            }
        }
    }

    return Point{x: -1, y: -1} // Indicating asteroid not found!
}

func HCF(n1, n2 int) int {
    if n1 < 0 { n1 = -n1 }
    if n2 < 0 { n2 = -n2 }

    for n2 != 0 {
        n1, n2 = n2, n1 % n2
    }

    return n1
}

