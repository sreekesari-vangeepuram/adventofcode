package main

import "fmt"

func main() {
	// Password range -> [ 183564 ... 657474 ]

	// Rules:
	/*  1. Six digits long. [No need to check this, since the given range only consists of 6 digited numbers]
	 *  2. Within the above mentioned range. [Of course!]
	 *  3. Two adjacent digits are same (like `22` in `122345`).    [So this is our R1]
	 *  4. Going from left to right, the digits never decrease; they only ever
	 *     increase or stay the same (like `111123` or `135679`).   [Another rule R2]
	 */

	var groupSize, previousDigit, bufferNum, acc1, acc2 int // max size of both accumulators in the given range is `473910`
	var inIncreasingOrder, hasSimilarDigit, validGroup bool
	for number := 183564 + 1; number < 657474; number++ {
		bufferNum, previousDigit, groupSize = number, -1, 0
		inIncreasingOrder, hasSimilarDigit, validGroup = true, false, false
		for i := 5; i >= 0; i-- {
			previousDigit = bufferNum % 10
			bufferNum /= 10

			if previousDigit < bufferNum%10 {
				inIncreasingOrder, validGroup = false, false
				break
			}

			if previousDigit == bufferNum%10 {
				hasSimilarDigit = true
				groupSize++

			} else {
				if groupSize == 1 {
					validGroup = true
				}
				groupSize = 0
			}

		}

		if inIncreasingOrder && hasSimilarDigit {
			acc1 += 1
			if groupSize == 1 || validGroup {
				acc2 += 1
			}
		}
	}

	fmt.Printf("Number of different passwords obeying the rules: %d\n", acc1)
	fmt.Printf("Number of different passwords obeying the modified-rules: %d\n", acc2)
}
