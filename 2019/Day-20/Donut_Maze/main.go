package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Vec2d struct { x, y int }

func (v Vec2d) Scale(v0 Vec2d) Vec2d {
	return Vec2d{
		x: v.x + v0.x,
		y: v.y + v0.y,
	}
}

var directions []Vec2d = []Vec2d{
                   /*UP*/
                Vec2d{0,  1},
/*LEFT*/ Vec2d{-1, 0}, Vec2d{1, 0}, /*RIGHT*/
                Vec2d{0, -1},
                   /*DOWN*/
}

type Node struct {
	location	Vec2d
	links		[]Link
}

type Link struct {
	neighbour	*Node
	offset		int
}

type Label struct {
	id			string
	nodes		[]*Node
}

type Location struct {
	node		*Node
	level		int
}

type Cell struct {
	location	Location
	distance	int
}

func main() {

	if len(os.Args) < 2 {
		log.Fatal(`
[ERROR]: Provide the input dataset!
**Usage: ./main /path/to/file
`)
	}

	rows := ReadFile(os.Args[1])
	initialNode, finalNode := locateAndLink(rows)

	fmt.Println("Number of steps from AA to ZZ of open tile marked :", calculateDistance(initialNode, finalNode, false))
	fmt.Println("Number of steps from AA to ZZ of outermost-layers :", calculateDistance(initialNode, finalNode,  true))
}

func locateAndLink(rows []string) (*Node, *Node) {
	WIDTH, DEPTH := len(rows[0]), len(rows)

	nodes               := make(map[Vec2d]*Node)
	labelsPositionMap   := make(map[Vec2d]*Label)
	labelsIdentifierMap := make(map[string]*Label)

	var label *Label

	// Labelling and identification of all (useful) nodes
	for y, row := range rows {
		for x, block := range row {

			if block != '.' { continue }

			loc        := Vec2d{x, y}
			node       := &Node{location: loc}
			nodes[loc]  = node

			for _, dir := range directions {

				next      := loc.Scale(dir)
				nextBlock := rows[next.y][next.x]

				if !(nextBlock >= 'A' && nextBlock <= 'Z') { continue }

				tail      := next.Scale(dir)
				tailBlock := rows[tail.y][tail.x]

				id := string(tailBlock) + string(nextBlock)
				if dir.y >= 0 && dir.x >= 0 {
					id = string(nextBlock) + string(tailBlock)
				}

				label = labelsIdentifierMap[id]
				if label == nil {
						label = &Label{id: id}
						labelsIdentifierMap[id] = label
				}

				labelsPositionMap[next] = label
				label.nodes = append(label.nodes, node)
			}
		}
	}

	// Linking nodes
	for _, node := range nodes {
		for _, dir := range directions {

			nextLocation := node.location.Scale(dir)
			nextNode     := nodes[nextLocation]
			nextLabel    := labelsPositionMap[nextLocation]

			if nextNode  != nil { node.links = append(node.links, Link{nextNode, 0}) }
			if nextLabel != nil {
				offset := 1
				if nextLocation.x <= 2 || nextLocation.y <= 2 { offset = -1; goto LOOP }
				if nextLocation.x >= WIDTH-2 || nextLocation.y >= DEPTH-2 { offset = -1 }

				LOOP: for _, nextNode := range nextLabel.nodes {
					if nextNode != node {
						node.links = append(node.links, Link{
								neighbour : nextNode,
								offset    : offset,
						});		break LOOP
					}
				}
			}
		}
	}

	return labelsIdentifierMap["AA"].nodes[0], labelsIdentifierMap["ZZ"].nodes[0]
}

func calculateDistance(initialNode, finalNode *Node, rec bool) int {

	initialCell := Cell{
		location: Location{node : initialNode, level: 0},
		distance: 0,
	}

	var openCells []Cell
	openCells = append(openCells, initialCell)

	visited := make(map[Location]bool)
	visited[initialCell.location] = true

	for/*ever*/ {

		cell      := openCells[0]
		openCells  = openCells[1:]

		if cell.location.node == finalNode && cell.location.level == 0 { return cell.distance }

		for _, link := range cell.location.node.links {

			nextLocation := Location{node : link.neighbour, level: cell.location.level}
			if rec { nextLocation.level += link.offset }

			if nextLocation.level >= 0 && !visited[nextLocation] {
				openCells = append(openCells, Cell{location: nextLocation, distance: cell.distance + 1})
				visited[nextLocation] = true
			}
		}
	}

	return -1 // In-case of invalid input!
}

func ReadFile(fileName string) (rows []string) {
	fhand, err := os.Open(fileName)
	if err != nil { log.Fatal(err) }
	defer fhand.Close()

	scanner := bufio.NewScanner(fhand)
	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}; return rows
}
