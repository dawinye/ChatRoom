package main

import (
	"bufio"
	"strings"

	//"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
)

//message struct is defined in both the server and client files
//this is because it is classified as undefined if it is only defined in one
//and we use either go run server.go or go run client.go
//(error pops up for the file that doesn't have it defined)
type Message struct {
	To, From, Content string
}

//keep a map of usernames as keys and their connections as values
//by using a map, we can check if the recipient of a message is already
//defined in a server, as well as check if a username should be rejected
//if it already exists in the server
var m = make(map[string]net.Conn)
var finish bool = false

//use a connection handler to ensure each server can handle multiple clients
func handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	text, _ := reader.ReadString('\n')
	username := strings.TrimSpace(string(text))
	m[username] = conn

	//defer conn.Close()
	tmp := make([]byte, 500)
	for {
		_, err := conn.Read(tmp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tmpbuff := bytes.NewBuffer(tmp)
		tmpstruct := new(Message)
		gobobj := gob.NewDecoder(tmpbuff)
		gobobj.Decode(tmpstruct)
		sendMessage(*tmpstruct)
	}
}

//input is the message that the client sends back to the server
//server checks if the recipient of the message is in the map
//sends an error message to the sender if applicable, otherwise
//sends a success message and delivers it to the recipient connection
func sendMessage(msg Message) {
	fromConn := m[msg.From]
	toConn, present := m[msg.To]

	if present == false {
		fromConn.Write([]byte("Failure"))
	} else {
		toConn.Write([]byte("From: " + msg.From + " Message: " + msg.Content))
		fromConn.Write([]byte("Message successfully delivered to " + msg.To))
	}
}
func stopServer() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(string(text)) == "EXIT" {
			closeClients()
			fmt.Println("Server terminated.")
			os.Exit(1)
			return
		} else {
			fmt.Println("Please enter \"EXIT\" to stop the server. No other commmands are supported")
		}
	}
}
func closeClients() {
	for _, conn := range m {
		conn.Write([]byte("EXIT"))
		fmt.Println([]byte("EXIT"))
		fmt.Println(string([]byte("EXIT")))
	}
}

func main() {
	args := os.Args

	//error checking to see if the port number is provided
	if len(args) != 2 {
		fmt.Println("Please rerun the program using \"go run server.go (port number)\"")
		return
	}

	//check if the port number supplied falls within the preallocated ports
	port, err := strconv.Atoi(args[1])
	if err != nil || port < 0 || port > 65535 {
		fmt.Println("Please rerun the program using \"go run server.go (port number between 0 and 65535)\"")
		return
	}
	l, err := net.Listen("tcp4", ":"+args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	go stopServer()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		//run each client connection in a goroutine to allow concurrency
		go handleConn(conn)
	}

}
