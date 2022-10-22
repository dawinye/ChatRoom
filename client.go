package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strings"
)

type Message struct {
	To, From, Content string
}

func main() {
	//check the user inputs to see if they entered the right amount
	args := os.Args
	if len(args) != 4 {
		fmt.Println("Please rerun the program using \"go run client.go (host address) (port number) (username)")
		return
	}

	//assign the host, port and username
	host := args[1]
	port := args[2]
	username := args[3]

	//limit host addresses to only the local machine
	if host != "127.0.0.1" {
		fmt.Println("Please use 127.0.0.1 as the host to connect to the local server")
		return
	}

	c, err := net.Dial("tcp4", host+":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	//tmp := make([]byte, 500)
	for {
		//bin_buf := new(bytes.Buffer)

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(string(text)) == "EXIT" {
			fmt.Println("TCP client exiting...")
			return
		}
		messageArr := strings.Split(text, ",,")
		if len(messageArr) != 2 {
			// example: "Alan,,you are awful at poker" sends "you are awful at poker" to Alan if he is connected to the server
			fmt.Println("Please enter your message with two commas (,,) to separate the recipient of the message and the message contents  (in that order)")
			continue
		}
		msg := new(Message)
		msg.To = messageArr[0]
		msg.Content = messageArr[1]
		msg.From = username
		bin_buf := new(bytes.Buffer)
		gobobj := gob.NewEncoder(bin_buf)
		gobobj.Encode(msg)

		c.Write(bin_buf.Bytes())

		//fmt.Fprintf(c, msg.content+"\n")
		//message, _ := bufio.NewReader(c).ReadString('\n')
		//fmt.Print("->: " + message)

	}
	//for {
	//	reader := bufio.NewReader(os.Stdin)
	//	fmt.Print(">> ")
	//	text, _ := reader.ReadString('\n')
	//	fmt.Fprintf(c, text+"\n")
	//
	//	message, _ := bufio.NewReader(c).ReadString('\n')
	//	fmt.Print("->: " + message)
	//	if strings.TrimSpace(string(text)) == "STOP" {
	//		fmt.Println("TCP client exiting...")
	//		return
	//	}
	//}
}
