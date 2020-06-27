package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {

	// go run udpclient.go 127.0.0.1:1234
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	CONNECT := arguments[1]

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	c, err := net.DialUDP("udp4", nil, s)

	//c, err := net.Dial("udp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("udp client start success")
	fmt.Printf("The UDP client is %s\n", c.RemoteAddr().String())
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">>")
		text, _ := reader.ReadString('\n')
		data := []byte(text + "\n")
		_, err := c.Write(data)

		if strings.TrimSpace(string(data)) == "STOP" {
			fmt.Println("Exiting UDP CLIENT!")
			return
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		// new 的作用是初始化一个指向类型的指针(*T)，
		// make 的作用是为 slice，map 或 chan 初始化并返回引用(T)
		buffer := make([]byte, 1024)
		n, _, err := c.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Reply: %s\n", string(buffer[0:n]))

	}

}
