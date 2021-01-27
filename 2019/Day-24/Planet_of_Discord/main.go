package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Point struct{ x, y int }

func (p Point) Add(p0 Point) Point {
	return Point{
		x: p.x + p0.x,
		y: p.y + p0.y,
	}
}

var directions []Point = []Point{
	                   /*UP*/
	                Point{0, 1},
	/*LEFT*/ Point{-1, 0}, Point{1, 0}, /*RIGHT*/
	                Point{0, -1},
	                  /*DOWN*/
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal(`
[ERROR]: Provide the input dataset!
**Usage: ./main /path/to/file
`)
	}

	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil { log.Fatal(err) }

	Eris := make(map[Point]int)
	for y, row := range strings.Fields(string(bytes)) {
		for x, block := range row {
			switch block {
			case '.':
				Eris[Point{x, y}] = 0

			case '#':
				Eris[Point{x, y}] = 1

			default:
				log.Fatal("[ERROR]: Unknown block encountered from the scan of Eris!")
			}
		}
	}

	var PlutonianGrid [5][5]int
	for y, row := range strings.Fields(string(bytes)) {
		for x, block := range row {
			if block /*has bug*/ == '#' { PlutonianGrid[y][x] = 1 }
		}
	}

	fmt.Println("Biodiversity rating for the first layout that appears twice :", GameOfLife(Eris, 5, 5))
	fmt.Println("Number of bugs present after 200 minutes :", initScan(PlutonianGrid, 200 /*minutes*/, 5, 5))

}

func GameOfLife(Eris map[Point]int, WIDTH, DEPTH int) (rating int) {

	scanned  := make(map[int]interface{})
	Infected := make(map[Point]int)
	var adjacentBlock int

	for /*ever*/ {

		rating = 0
		Infected = make(map[Point]int)

		for point, rate := range Eris {

			adjacentBlock = 0
			for _, dir := range directions {
				adjacentBlock += Eris[point.Add(dir)]
			}

			Infected[point] = 0
			if adjacentBlock == 2 && rate == 0 || adjacentBlock == 1 {
				Infected[point] = 1
				rating += 1 << ((DEPTH * point.y) + point.x)
			}
		}; Eris = Infected
		_, done := scanned[rating]

		if done {
			return rating
		}
		scanned[rating] = nil
	}

	return 0 // No Rate...
}

func initScan(PlutonianGrid [5][5]int, time int, WIDTH, DEPTH int) (BugsCount int) {

	Grid     := make(map[int][5][5]int)
	Infected := make(map[int][5][5]int)

	var adjacentBlock, X, Y int
	var p Point

	Grid[0] = PlutonianGrid // Since it is recursive...
	for i := 0; i < time; i++ {

		BugsCount = 0
		Infected = make(map[int][5][5]int)
		var bufferGrid [5][5]int

		for t := -time; t <= time; t++ {
			for y := 0; y < DEPTH; y++ {
				for x := 0; x < WIDTH; x++ {

					if x == 2 && y == 2 { continue } // Sub-grid marked as -> `?`

					adjacentBlock = 0
					for _, dir := range directions {

						p = Point{x, y}.Add(dir)

						if p.x == -1 || p.y == -1 || p.x == 5 || p.y == 5 {
							adjacentBlock += Grid[t - 1][2 + dir.y][2 + dir.x]
						} else if p.x == 2 && p.y == 2 {
							for j := 0; j < DEPTH; j++ {

								X = dir.y & 1*j - dir.x & 1*2*(dir.x - 1)
								Y = dir.x & 1*j - dir.y & 1*2*(dir.y - 1)

								adjacentBlock += Grid[t + 1][Y][X]
							}
						} else {adjacentBlock += Grid[t][p.y][p.x]}
					}

					if adjacentBlock == 2 && Grid[t][y][x] == 0 || adjacentBlock == 1 {
						bufferGrid = Infected[t]
						bufferGrid[y][x] = 1; Infected[t] = bufferGrid
						BugsCount++
					}
				}
			}
		}

	    Grid = Infected // PreviousOuput -> PresentInput ... Recursive!
	};  return BugsCount
}
