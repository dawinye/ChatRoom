package main

import (
	"bufio"
	"io"
	"strings"

	//"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
)

// message struct is defined in both the server and client files
// this is because it is classified as undefined if it is only defined in one
// and we use either go run server.go or go run client.go
// (error pops up for the file that doesn't have it defined)
type Message struct {
	To, From, Content string
}

// keep a map of usernames as keys and their connections as values
// by using a map, we can check if the recipient of a message is already
// defined in a server, as well as check if a username should be rejected
// if it already exists in the server
var m = make(map[string]net.Conn)

// use a connection handler to ensure each server can handle multiple clients
func handleConn(conn net.Conn) {
	//made a buffer and store what's read from the client into said buffer
	usernameBuf := make([]byte, 500)
	conn.Read(usernameBuf)

	//decode the username and populate the map
	username := new(string)
	tmpBuff := bytes.NewBuffer(usernameBuf)
	gobObj := gob.NewDecoder(tmpBuff)
	gobObj.Decode(username)
	m[*username] = conn

	fmt.Println(*username + " has connected to this server")
	fmt.Print(">> ")

	//store new messages into tmp
	tmp := make([]byte, 500)
	for {
		_, err := conn.Read(tmp)

		// if the client disconnects unexpectedly via ctrl + C or voluntarily via exit, it disconnects
		//remove the mapping and close the connection
		if err == io.EOF || string(tmp[:4]) == "EXIT" {
			delete(m, *username)
			conn.Close()
			fmt.Println(*username + " has disconnected from this server")
			fmt.Print(">> ")
			return
		}
		if err != nil {
			fmt.Println(err)
			continue
		}

		//decode the message from the client
		tmpBuff = bytes.NewBuffer(tmp)
		tmpStruct := new(Message)
		gobObj = gob.NewDecoder(tmpBuff)
		gobObj.Decode(tmpStruct)
		//send it to this function which checks if user2 is connected
		sendMessage(*tmpStruct)

	}
}

// input is the message that the client sends back to the server
// server checks if the recipient of the message is in the map
// sends an error message to the sender if applicable, otherwise
// sends a success message and delivers it to the recipient connection
func sendMessage(msg Message) {
	fromConn := m[msg.From]
	toConn, present := m[msg.To]

	//separate encoders and buffers
	bin_buf := new(bytes.Buffer)
	bin_bufTwo := new(bytes.Buffer)
	gobobj := gob.NewEncoder(bin_buf)
	gobobjTwo := gob.NewEncoder(bin_bufTwo)

	if present == false {
		strOne := "Error: That user is not connected to the server. Maybe they will be here soon!"
		gobobj.Encode(strOne)
		fromConn.Write(bin_buf.Bytes())
	} else if msg.From == msg.To {
		strOne := "Cannot send messages to yourself. Use a notepad to take notes instead of this program."
		gobobj.Encode(strOne)
		fromConn.Write(bin_buf.Bytes())
	} else {
		strOne := "Message successfully delivered to " + msg.To
		gobobj.Encode(strOne)
		fromConn.Write(bin_buf.Bytes())

		strTwo := "From: " + msg.From + " Message: " + msg.Content
		gobobjTwo.Encode(strTwo)
		toConn.Write(bin_bufTwo.Bytes())
	}
}

//function to check if the user enters EXIT
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

//closes all client connection if the server terminates
func closeClients() {
	for _, conn := range m {
		conn.Write([]byte("EXIT"))
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
	if err != nil || port < 1 || port > 65535 {
		fmt.Println("Please rerun the program using \"go run server.go (port number between 1 and 65535, inclusive)\"")
		return
	}

	l, err := net.Listen("tcp4", ":"+args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	//check if the server stops in a goroutine
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
