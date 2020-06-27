package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const MIN = 1
const MAX = 100

func random() int {
	return rand.Intn(MAX-MIN) + MIN
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		temp := strings.TrimSpace(string(netData))
		log.Println(c.RemoteAddr().String() + " " + temp)
		if temp == "STOP" {
			break
		}

		result := strconv.Itoa(random()) + "\n"
		c.Write([]byte(string(result)))
	}
	c.Close()
}

func main() {

	// arguments := os.Args
	// if len(arguments) == 1 {
	// 		fmt.Println("Please provide a port number!")
	// 		return
	// }

	// PORT := ":" + arguments[1]
	PORT := ":8001"
	// net
	l, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(c)
	}
}
