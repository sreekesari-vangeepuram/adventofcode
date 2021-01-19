package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal(`
[ERROR]: Provide the input dataset!
**Usage: ./main /path/to/file
`)
	}

	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil { log.Fatal(err) }

	bytes = []byte(strings.Trim(string(bytes), " \n\r\t")) // Sanitation

	customizeBytes := func(bytes *[]byte) {
		for i, c := range *bytes {
			(*bytes)[i] = byte(c-'0')
		}
	}

	customizeBytes(&bytes)

	RealSignal := make([]byte, 0, len(bytes) * 10000)
	for i := 0; i < 10000; i++ {
		RealSignal = append(RealSignal, bytes...)
	}

	var buffString string = ""
	for i := 0; i < 7; i ++ {
		buffString += string(bytes[i]+'0')
	}

	offset := Int(buffString)

	fmt.Println("First eight digits in the final output list:", Phaser(bytes, 100)[:8])
	fmt.Println("Eight-digit message embedded in the final output list:", decodeRealSignal(RealSignal, offset, 100)[:8])

}

func decodeRealSignal(RealSignal []byte, offset, repeatitions int) []byte {
	// For efficiency's sake
	// skipping the element until the
	// offset boost the runtime
	RealSignal = RealSignal[offset:]
	for phase := 1; phase <= repeatitions; phase++ {
		acc := 0
		// Faster decoding from the RHS-end
		for i := len(RealSignal) - 1; i >= 0; i-- {
			acc += int(RealSignal[i])
			RealSignal[i] = byte(abs(acc % 10))
		}
	}

	return RealSignal
}

func Phaser(InputSignal []byte, repeatitions int) []byte {

	for phase := 1; phase <= repeatitions; phase++ {

		OutputSignal := make([]byte, len(InputSignal))
		for i := 0; i < len(InputSignal); i++ {

			acc := 0
			for j := i; j < len(InputSignal); {

				for p := 0; j < len(InputSignal) && p < i + 1; p++ {
					acc += int(InputSignal[j]); j++
				}; j += i + 1

				for p := 0; j < len(InputSignal) && p < i + 1; p++ {
					acc -= int(InputSignal[j]);	j++
				}; j += i + 1

			}; OutputSignal[i] = byte(abs(acc % 10)) // One's digit of the absolute-summation
		}
			InputSignal = OutputSignal				 // InputSignal for the next phase is OutputSignal
	}

	return InputSignal
}

func abs(n int) int {
	if n < 0 { return -n }
	return n
}

func Int(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil { log.Fatal(err) }
	return n
}
