package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Message struct {
	To, From, Content string
}

// use to help slice the decoded string from gob
func sliceHelper(sl []byte) []byte {
	low := 0
	high := len(sl) - 1
	mid := 0
	// 10 is Ascii newLine character, which signifies end of the message
	target := 10

	// binary search for target; receiver looks like: [{positive nums != 10}, 10, 0, 0, ... 0]
	for low <= high {
		mid = (high + low) / 2
		if int(sl[mid]) == target {
			return sl[:mid]
		}
		// a positive value (!= 10) means index mid is still indexing part of the message
		if int(sl[mid]) > 0 {
			low = mid + 1
			// 0 means index mid is not indexing the message
		} else if int(sl[mid]) == 0 {
			high = mid - 1
		}
	}
	// 10 not found; must be sender
	return sl
}

// reads from the connection and prints it out if another client sends a message
func receiveMessage(conn net.Conn) {
	defer conn.Close()
	msg := make([]byte, 500)
	//str := new(string)

	for {
		_, err := conn.Read(msg)
		if err == io.EOF {
			fmt.Println("Server has unexpectedly terminated the connection. Exiting now.")
			os.Exit(1)
		}
		if err != nil {
			fmt.Println("this is err: ", err)
		}
		tmpbuff := bytes.NewBuffer(msg)
		gobobj := gob.NewDecoder(tmpbuff)
		gobobj.Decode(msg)

		fmt.Println(string(sliceHelper(msg[4:])))
		fmt.Print(">> ")

		//clear the message array so that the next time we receive a message we don't overwrite the old one,
		//which could be problematic if the new message is shorter than the old one
		msg = make([]byte, 500)
	}
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

	//use a goroutine to continuously check if a message has been sent to this client
	go receiveMessage(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	bin_buf := new(bytes.Buffer)
	gobobj := gob.NewEncoder(bin_buf)
	gobobj.Encode(username)
	c.Write(bin_buf.Bytes())
	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(string(text)) == "EXIT" {
			c.Write([]byte("EXIT"))
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

		//populate the message struct with appropriate values
		msg.To = messageArr[0]
		msg.Content = messageArr[1]
		msg.From = username

		//use gob to encode the message and write to the server
		bin_buf = new(bytes.Buffer)
		gobobj = gob.NewEncoder(bin_buf)
		gobobj.Encode(msg)
		c.Write(bin_buf.Bytes())
	}

}
