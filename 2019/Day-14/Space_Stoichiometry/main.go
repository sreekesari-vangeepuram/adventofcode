package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "regexp"
    "strconv"
    "strings"
)

type Chemical struct {
    quantity    int
    name        string
}

type Reaction struct {
    reactants       []Chemical
    product         Chemical
}

func sliceToReaction(slice [][]string) Reaction {
    var reaction Reaction
    reaction.reactants = make([]Chemical, len(slice) - 1)

    for i, node := range slice {
        node = strings.Split(node[0], " ")
        n, err := strconv.Atoi(node[0])
        if err != nil { log.Fatal(err) }

        if i != len(slice) - 1 {
            reaction.reactants[i] = Chemical{
                quantity: n,
                name: node[1],
            }
        } else {
            reaction.product = Chemical{
                quantity: n,
                name: node[1],
            }
        }
    }

    return reaction
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

    data := strings.Split(strings.Trim(string(bytes), " \n"), "\n")
    pattern := regexp.MustCompile(`\d+ \w+`)

    var reaction Reaction
    reactions := make(map[string]Reaction)

    for _, line := range data {
        slice := pattern.FindAllStringSubmatch(line, -1)
        reaction = sliceToReaction(slice)
        reactions[reaction.product.name] = reaction
    }

    requiredMaterial := map[string]int{"FUEL": 1}
    decideOREquantity(reactions, &requiredMaterial)
    var OREsPer1FUEL int = requiredMaterial["ORE"]

    fmt.Println("Number of ORE-material required to produce 1 FUEL:", OREsPer1FUEL)
    fmt.Println("Maximum amount of FUEL that can be produced from 1 trillion ORE:", fetchMaxFuel(reactions, OREsPer1FUEL))
}

func decideOREquantity(reactions map[string]Reaction, requiredMaterial *map[string]int) {

    var reacted bool

    for {

        reacted = false
        for chemical, quantity := range *requiredMaterial {
            if quantity > 0 {
                reaction, present := reactions[chemical]
                if present {
                    reacted = true
                    frameScale := (quantity - 1)/reaction.product.quantity + 1
                    (*requiredMaterial)[chemical] -= reaction.product.quantity * frameScale
                    for _, reactant := range reaction.reactants {
                        (*requiredMaterial)[reactant.name] += reactant.quantity * frameScale
                    }
                }
            }
        }

        if !reacted { break }
    }
}

func fetchMaxFuel(reactions map[string]Reaction, OREsPer1FUEL int) (maxFUEL int) {
    chunkSize := 1000000000000/OREsPer1FUEL
    requiredMaterial := make(map[string]int)

    LOOP: for {
        copyOfRequiredMaterial := make(map[string]int)
        for chemical, quantity := range requiredMaterial { copyOfRequiredMaterial[chemical] = quantity }

        copyOfRequiredMaterial["FUEL"] += chunkSize
        decideOREquantity(reactions, &copyOfRequiredMaterial)

        if copyOfRequiredMaterial["ORE"] <= 1000000000000 {
            maxFUEL += chunkSize
            requiredMaterial = copyOfRequiredMaterial
            goto LOOP
        }

        if chunkSize > 1 { chunkSize /= 2; goto LOOP }

        return // ORE-material is totally used for the production of FUEL!
    }
}
