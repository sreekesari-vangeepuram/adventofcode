package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(`

[ERROR]: Provide input dataset!
** Usage: ./main /path/to/file
        `)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	bufferString := strings.TrimSpace(string(data))
	offset := 25 * 6

	imageData := make([]string, len(bufferString)/offset)
	for i := 0; len(bufferString) != 0; i++ {
		imageData[i] = bufferString[:offset]
		bufferString = bufferString[offset:]
	}

	checkSumString := imageData[findLeastCountLayer(imageData, "0")]
	oneDigits, twoDigits := 0, 0
	for _, s := range checkSumString {
		if string(s) == "1" {
			oneDigits++
		} else if string(s) == "2" {
			twoDigits++
		}
	}

	fmt.Println("Checksum:", oneDigits*twoDigits)
	fmt.Println("Image:")
	printImage(data)
}

func findLeastCountLayer(imageData []string, char string) (index int) {
	leastCount := 25 * 6
	for i, layer := range imageData {
		zeroCount := 0
		for j := 0; j < 25*6; j++ {
			if string(layer[j]) == char {
				zeroCount++
			}
		}
		if leastCount > zeroCount {
			leastCount = zeroCount
			index = i
		}
	}

	return index
}

func printImage(data []byte) {
	imageData := []byte(strings.Trim(string(data), "\n "))

	for i := range imageData {
		imageData[i] -= 48 // byte("0") = 48
	}

	image := make([]string, 25*6)
	for i, _ := range image {
		image[i] = "2"
	}

	for layer := 0; layer < len(imageData)/(25*6); layer++ {
		for row := 0; row < 6; row++ {
			for column := 0; column < 25; column++ {

				pixel := imageData[layer*(25*6)+row*25+column]
				if image[row*25+column] == "2" {
					if pixel == 0 {
						image[row*25+column] = " "
					} else if pixel == 1 {
						image[row*25+column] = "@"
					}
				}
			}
		}
	}

	for i := 1; i <= 25*6; i++ {
		fmt.Printf("%s", image[i-1])
		if i%25 == 0 {
			fmt.Println()
		}
	}

}
