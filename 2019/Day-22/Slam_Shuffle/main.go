package main

import (
	"bufio"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
)

type Instruction struct { scale, shift *big.Int }

func (ins Instruction) mod(scalar int64) {
	ins.scale.Mod(ins.scale, big.NewInt(scalar))
	ins.shift.Mod(ins.shift, big.NewInt(scalar))
}

func (ins Instruction) instruct(n1, n2 *big.Int) {
	n1.Mul(n1, ins.scale)
	n1.Add(n1, ins.shift)
	n1.Mod(n1, n2)
}


func (ins Instruction) invert(scalar int64) Negation {

	ins.mod(scalar)
	n, scale := big.NewInt(scalar), big.NewInt(0)
	if ins.scale.Cmp(scale) != 0 {
		scale = ins.scale.ModInverse(ins.scale, n)
	}

	shift := ins.shift.Mul(ins.shift, big.NewInt(-1))
	shift.Mul(shift, scale)
	shift.Mod(shift, n)

	return Negation{scale, shift}
}

type Negation Instruction // Inverted Instruction

func (neg Negation) instruct(n1, n2 *big.Int) {
	n1.Mul(n1, neg.scale)
	n1.Add(n1, neg.shift)
	n1.Mod(n1, n2)

}

func (neg Negation) pow(n1, n2 *big.Int) {
	var scale big.Int
	scale.Set(neg.scale)
	SUM(neg.shift, &scale, n1, n2)
	neg.scale.Exp(neg.scale, n1, n2)
}

func SUM(n1, n2, n3, n4 *big.Int) {

	var buff big.Int
	buff.Set(big.NewInt(1))
	buff.Sub(&buff, n2)
	buff.ModInverse(&buff, n4)

	n2.Exp(n2, n3, n4)
	n2.Mul(n2, big.NewInt(-1))
	n2.Add(n2, big.NewInt(1))

	n2.Mul(n2, &buff)
	n2.Mul(n2, n1)

	n1.Set(n2)
	n1.Mod(n1, n4)
}

type Instructor interface {
	instruct(n1, n2 *big.Int)	// function signature
}

func main() {

	if len(os.Args) < 2 {
log.Fatal(`
[ERROR]: Provide the input dataset!
**Usage: ./main /path/to/file
`)
	}

	instruction := ReadAndParseFile(os.Args[1])

	var (
		n1 int64 = 2019
		n2 int64 = 10007
		n  int   = 1
	)

	fmt.Printf("Card number at %d: %d\n", n1, Solve_1(n1, n2, instruction, n))


		n1 = 2020
		n2 = 119315717514047
		n  = 101741582076661

	fmt.Printf("Card number at %d: %v\n", n1, Solve_2(n1, n2, instruction, n))
}

func Solve_1(n1, n2 int64, intr Instructor, n int) int64 {

	i1, i2 := big.NewInt(n1), big.NewInt(n2)
	for ; n > 0; n-- {
		intr.instruct(i1, i2)
		i1.Mod(i1, i2)
	}

	return i1.Int64()
}

func Solve_2(n1, n2 int64, ins Instruction, n int) int64 {

	i1  := ins.invert(n2)
	i2  := big.NewInt(n2)
	neg := big.NewInt(n1)

	i1.pow(big.NewInt(int64(n)), i2)
	i1.instruct(neg, i2)

	neg.Mod(neg, i2)
	return neg.Int64()
}

func ReadAndParseFile(fileName string) Instruction {
	fhand, err := os.Open(fileName)
	if err != nil { log.Fatal(err) }
	defer fhand.Close()

	scanner := bufio.NewScanner(fhand)
	var line string
	var buffBigNum big.Int
	var ins Instruction = Instruction{
			scale: big.NewInt(1),
			shift: big.NewInt(0),
	}

	for ; scanner.Scan();  {

		line = scanner.Text()
		var buffNum int64 = 0

		if strings.Index(line, "deal into") == 0 {
			ins.scale.Mul(ins.scale, big.NewInt(-1))
			ins.shift.Mul(ins.shift, big.NewInt(-1))
			ins.shift.Sub(ins.shift, big.NewInt(1))
		} else if strings.Index(line, "cut") == 0 {
			fmt.Sscanf(line, "cut %d", &buffNum)
			buffBigNum.SetInt64(buffNum)
			ins.shift.Sub(ins.shift, &buffBigNum)
		} else if strings.Index(line, "deal with") == 0 {
			fmt.Sscanf(line, "deal with increment %d", &buffNum)
			buffBigNum.SetInt64(buffNum)
			ins.scale.Mul(ins.scale, &buffBigNum)
			ins.shift.Mul(ins.shift, &buffBigNum)
		} else { log.Fatal("[ERROR]: Unknown instruction!") }
	}

	if err := scanner.Err(); err != nil { log.Fatal(err) }
	return ins
}
