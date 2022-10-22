package main

import (
	//"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Message struct {
	To, From, Content string
}

//func handleConn(conn net.Conn) {
//	defer conn.Close()
//
//	for {
//
//	}
//}

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("Please rerun the program using \"go run server.go (port number)\"")
		return
	}
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

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	tmp := make([]byte, 500)
	for {
		_, err = c.Read(tmp)
		tmpbuff := bytes.NewBuffer(tmp)
		tmpstruct := new(Message)
		gobobj := gob.NewDecoder(tmpbuff)
		gobobj.Decode(tmpstruct)
		fmt.Println(tmpstruct)

		//netData, err := bufio.NewReader(c).ReadString('\n')
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//fmt.Println(strings.TrimSpace(string(netData)))
		//if strings.TrimSpace(string(netData)) == "STOP" {
		//	fmt.Println("Exiting TCP server!")
		//	return
		//}
		//
		//fmt.Print("-> ", string(netData))
		//t := time.Now()
		//myTime := t.Format(time.RFC3339) + "\n"
		//c.Write([]byte(myTime))
	}
}
