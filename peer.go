package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

type Message struct {
	Data string
}

func main() {
	finished := make(chan bool)

	args := os.Args[1:]
	port := args[0]

	// for purpose of verbosity, I will be removing error handling from this
	// sample code

	// Printer session.
	go func() {
		server, _ := net.Listen("tcp", ":"+port)

		for {
			conn, _ := server.Accept()
			// create a temp buffer
			tmp := make([]byte, 500)
			_, _ = conn.Read(tmp)
			// convert bytes into Buffer (which implements io.Reader/io.Writer)
			tmpbuff := bytes.NewBuffer(tmp)
			tmpstruct := new(Message)

			// creates a decoder object
			gobobj := gob.NewDecoder(tmpbuff)

			// decodes buffer and unmarshals it into a Message struct
			gobobj.Decode(tmpstruct)

			// lets print out!
			fmt.Println("\nReceived message!")
			fmt.Println(tmpstruct) // reflects.TypeOf(tmpstruct) == Message{}
		}
	}()

	go func() {
		for {
			var port string
			var data string

			fmt.Print("Enter peer port: ")
			fmt.Scanf("%s", &port)
			fmt.Println("Port is", port)

			fmt.Print("Enter data: ")
			fmt.Scanf("%s", &data)
			fmt.Println("Data is", data)

			// connect to server
			conn, _ := net.Dial("tcp", "127.0.0.1:"+port)

			msg := Message{Data: data}

			bin_buf := new(bytes.Buffer)
			// create a encoder object
			gobobj := gob.NewEncoder(bin_buf)

			// encode buffer and marshal it into a gob object
			gobobj.Encode(msg)

			conn.Write(bin_buf.Bytes())

			conn.Close()
		}
	}()

	<-finished
}
