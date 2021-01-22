package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

type Vec2d struct{ x, y int }

func (v Vec2d) Add(v0 Vec2d) Vec2d {
	return Vec2d{
		x: v.x + v0.x,
		y: v.y + v0.y,
	}
}

var directions []Vec2d = []Vec2d{
                   /*UP*/
                Vec2d{0, 1},
/*LEFT*/ Vec2d{-1, 0}, Vec2d{1, 0}, /*RIGHT*/
                Vec2d{0, -1},
                   /*DOWN*/
}

type Block struct {
	char     byte
	location Vec2d
	path     []Track // Path to connect two blocks
}

type Track struct {
	key      *Block
	distance int
	l2c      uint32 // Length to connect
}

func checkBlock(block byte, bType string) bool {

	switch bType {
	case "entrance":
		return block == '@' || block == '|' || block == '&' || block == '^'

	case "key":
		return block >= 'a' && block <= 'z'

	case "door":
		return block >= 'A' && block <= 'Z'

	}

	return false
}

func distFrom(bType string, block byte) uint32 {

	switch bType {
	case "key":
		return uint32(1 << (block - 'a'))

	case "door":
		return uint32(1 << (block - 'A'))
	}

	return uint32(0)
}

func fetchEntraceLoction(vault [][]byte) Vec2d {
	for y, row := range vault {
		for x, block := range row {
			if block == '@' { return Vec2d{x, y} }
		}
	}

	return Vec2d{-1, -1} // No entrance found
}

func check6sides(vault [][]byte, entrance Vec2d) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 { continue } // Intial entrance location
			if vault[entrance.y+dy][entrance.x+dx] != '.' {
				log.Fatal("[ERROR]: Unable to send robot inside entrance at: %v\n", entrance.Add(Vec2d{dx, dy}))
			}
		}
	}
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal(`
[ERROR]: Please provide the input dataset!
**Usage: ./main /path/to/file
`)
	}

	parsedData := ReadFile(os.Args[1])

	vault := make([][]byte, len(parsedData))
	for y, row := range parsedData { vault[y] = []byte(row)	}

	fmt.Println("Length of shortest path that collects all of the keys [in vault steps]:", navigateVault(vault))

	E := fetchEntraceLoction(vault)

	check6sides(vault, E) // Some verification
	vault[E.y-1][E.x-1] = '@'; vault[E.y-1][E.x] = '#'; vault[E.y-1][E.x+1] = '|';
	vault[ E.y ][E.x-1] = '#'; vault[ E.y ][E.x] = '#'; vault[ E.y ][E.x+1] = '#';
	vault[E.y+1][E.x-1] = '&'; vault[E.y+1][E.x] = '#'; vault[E.y+1][E.x+1] = '^';

	fmt.Println("Number of fewest steps necessary to collect all of the keys [in steps]:", navigateVault(vault))

}

func navigateVault(vault [][]byte) int {

	var keysCount int = 0
	blocks := make(map[byte]*Block)

	// Aligning data from the grid
	for y, row := range vault {
		for x, block := range row {
			if checkBlock(block, "key") { keysCount++ }

			if checkBlock(block, "entrance") || checkBlock(block, "key") {
				blocks[block] = &Block{
					char:     block,
					location: Vec2d{x, y},
				}
			}
		}
	}

	type Item struct {
		location Vec2d
		l2c      uint32 // length top connect
		distance int
	}

	// Track and connect blocks through vault-exploration.
	for _, block := range blocks {

		var openBlocks []Item
		openBlocks = append(openBlocks, Item{block.location, 0, 0})

		visited := make(map[Vec2d]bool)
		visited[block.location] = true

		for len(openBlocks) != 0 {

			currentBlock := openBlocks[0]
			openBlocks = openBlocks[1:]

			for _, dirVec := range directions {

				nextVec := currentBlock.location.Add(dirVec)
				item := vault[nextVec.y][nextVec.x]
				var bufferItem Item

				if visited[nextVec] { continue }

				if item == '.' || checkBlock(item, "entrance") {
					bufferItem = Item{
						location: nextVec,
						l2c:      currentBlock.l2c,
						distance: currentBlock.distance + 1,
					}

					openBlocks = append(openBlocks, bufferItem)
					visited[nextVec] = true

				} else if checkBlock(item, "door") {
					bufferItem = Item{
						location: nextVec,
						l2c:      currentBlock.l2c | distFrom("door", item),
						distance: currentBlock.distance + 1,
					}

					openBlocks = append(openBlocks, bufferItem)
					visited[nextVec] = true

				} else if checkBlock(item, "key") {
					bufferItem = Item{
						location: nextVec,
						l2c:      currentBlock.l2c | distFrom("key", item),
						distance: currentBlock.distance + 1,
					}

					openBlocks = append(openBlocks, bufferItem)
					visited[nextVec] = true

					block.path = append(block.path, Track{
						key:      blocks[item],
						l2c:      currentBlock.l2c,
						distance: currentBlock.distance + 1,
					})

				}

			}

		}

	}

	var totalKeysCount uint32 = (1 << keysCount) - 1
	var locations []*Block

	for _, block := range blocks {
		if checkBlock(block.char, "entrance") { locations = append(locations, block) }
	}

	return UnlockDoors(locations, 0, math.MaxInt32, 0, totalKeysCount)
}

func UnlockDoors(locations []*Block, currentDistance, shortestDistance int, UnlockedDoorsCount, totalKeysCount uint32) int {

	IsUnlocked := func(doors, keys uint32) bool { return doors&keys == keys }

	if IsUnlocked(UnlockedDoorsCount, totalKeysCount) { return currentDistance }

	for i, loc := range locations {
		for _, unit := range loc.path {

			if !IsUnlocked(UnlockedDoorsCount, unit.l2c) { continue }
			if IsUnlocked(UnlockedDoorsCount, distFrom("key", unit.key.char)) {	continue }

			newDistance := currentDistance + unit.distance
			if shortestDistance < newDistance { continue }

			newLocations := make([]*Block, len(locations))
			copy(newLocations, locations)

			newLocations[i] = unit.key
			newUnlockedDoorsCount := UnlockedDoorsCount | distFrom("key", unit.key.char)

			finalDistance := UnlockDoors(newLocations, newDistance, shortestDistance, newUnlockedDoorsCount, totalKeysCount)
			if finalDistance < shortestDistance { shortestDistance = finalDistance }
		}
	}

	return shortestDistance
}

func ReadFile(fileName string) []string {

	fhand, err := os.Open(fileName)
	if err != nil { log.Fatal(err) }
	defer fhand.Close()

	var data []string
	scanner := bufio.NewScanner(fhand)

	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	return data

}
