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
[ERROR]: Provide the dataset!
**Usage: ./main /path/to/file
`)
	}

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	bufferArray := strings.Split(strings.Trim(string(data), " \n"), ",")
	Program := make([]int, len(bufferArray))
	for i, s := range bufferArray {
		n, _ := strconv.Atoi(s)
		Program[i] = n
	}

	// Initialize monitoring network traffic...
	monitorTraffic(Program, 255, 50)
}

func monitorTraffic(Program []int, MachineAddress, fleetSize int) {

	InboundStream := make([]chan int, fleetSize)
	OutboundStream := make([]chan int, fleetSize)
	SIGTERM := make(chan bool)

	// Inbound Packet Traffic
	for i := 0; i < fleetSize; i++ {

		InboundStream[i], OutboundStream[i] = make(chan int), make(chan int)
		go runProgram(Program, InboundStream[i], OutboundStream[i], SIGTERM)

		InboundStream[i] <- i
		InboundStream[i] <- -1
	}

	var NAT, oldAddress, newAddress [2]int /* Size = 2 for X & Y */
	var idle int = 0

	for i := 0; ; /* until receiving SIGTERM */ i = (i + 1) % fleetSize {
		select {
		case Address := <-OutboundStream[i]:
			if Address == MachineAddress {
				newAddress = [2]int{<-OutboundStream[i], <-OutboundStream[i]}
				if NAT == [2]int{} {
					fmt.Println("Y value of the first packet sent to address 255 :", newAddress[1])
				}
				NAT = newAddress
			} else {
				InboundStream[Address] <- int(<-OutboundStream[i])
				InboundStream[Address] <- int(<-OutboundStream[i])
			}
			idle = 0

		case InboundStream[i] <- -1:
			idle++

		case <-SIGTERM:
			return
		}

		if idle >= fleetSize {
			if oldAddress[1] == NAT[1] {
				fmt.Println("The first Y value delivered by the NAT to the computer at address 0 twice in a row :", NAT[1])
				return // DONE!
			}

			InboundStream[0] <- NAT[0]
			InboundStream[0] <- NAT[1]
			oldAddress = NAT
			idle = 0
		}
	}

	log.Fatal("[ERROR]: Unable to monitor traffic!")
}

func runProgram(Program []int, InputPacket <-chan int, OutputPacket chan<- int, SIGTERM chan<- bool) {
	mem := make([]int, len(Program)) // Expandable in `addr` func, on need only!
	copy(mem, Program)

	ip, rip := 0, 0         // instruction pointer [ip] and relative pointer [rip]
	var opCodeObject [5]int // array to split ABCDE
	var A, B, C, DE int     // digits: A, B, C in `ABCDE`
	var buffNum, p1, p2 int // buffer integer variables

	for {

		buffNum = int(mem[ip])
		for i := 4; i >= 0; i-- {
			opCodeObject[i] = buffNum % 10
			buffNum /= 10
		}

		DE = opCodeObject[3]*10 + opCodeObject[4] // opCode `DE` in `ABCDE` is extracted
		A, B, C = opCodeObject[0], opCodeObject[1], opCodeObject[2]

		p1, p2 = mem[addr(&mem, ip, C, 1, rip)], mem[addr(&mem, ip, B, 2, rip)]
		switch DE {

		case 1:
			mem[addr(&mem, ip, A, 3, rip)] = p1 + p2
			ip += 4

		case 2:
			mem[addr(&mem, ip, A, 3, rip)] = p1 * p2
			ip += 4

		case 3:
			mem[addr(&mem, ip, C, 1, rip)] = <-InputPacket
			ip += 2

		case 4:
			OutputPacket <- mem[addr(&mem, ip, C, 1, rip)]
			ip += 2

		case 5:
			if p1 != 0 {
				ip = p2
			} else {
				ip += 3
			}

		case 6:
			if p1 == 0 {
				ip = p2
			} else {
				ip += 3
			}

		case 7:
			if p1 < p2 {
				mem[addr(&mem, ip, A, 3, rip)] = 1
			} else {
				mem[addr(&mem, ip, A, 3, rip)] = 0
			}
			ip += 4

		case 8:
			if p1 == p2 {
				mem[addr(&mem, ip, A, 3, rip)] = 1
			} else {
				mem[addr(&mem, ip, A, 3, rip)] = 0
			}
			ip += 4

		case 9:
			rip += mem[addr(&mem, ip, C, 1, rip)]
			ip += 2

		case 99:
			SIGTERM <- true
			return

		default:
			log.Fatal("[ERROR]: Unknown opCode encountered!")

		}
	}

}

func addr(mem *[]int, ip int, digit int, offset int, rip int) (result int) {

	switch digit {
	case 0: // locationition mode
		result = (*mem)[ip+offset]

	case 1: // immediate mode
		result = ip + offset

	case 2: // relative mode
		result = (*mem)[ip+offset] + rip
	}

	// Expands the capacity on need
	for len(*mem) <= result {
		*mem = append(*mem, 0)
	}

	return abs(result)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
