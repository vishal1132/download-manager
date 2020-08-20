package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"os/user"

	"github.com/vishal979/downloadmanager/cmd/downloadmanager"
	"github.com/vishal979/downloadmanager/cmd/sockets"
)

//main is entry point of the program
func main() {
	var username string = "user"
	userProile, err := user.Current()
	if err == nil {
		username = userProile.Username
	}

	//make a socket connection
	//use this socket connection to send anything.
	conn, err := sockets.Run()
	if err != nil {
		conn = nil
	}
	var payload = make(map[string]interface{})

	func() {
		payload["Payload_type"] = "introduction"
		payload["Username"] = username
		payload["Body"] = "introduction payload"
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(&payload)
		conn.Write(buf.Bytes())
	}()

	if len(os.Args) == 1 {
		fmt.Println("Hi ", username, ", the proper use of the binary is <binary> <url> <optional:filelocation>")
		fmt.Println("You can also report or submit a review for this binary by <binary> <report/review> <your report/review>")
		return
	}
	if os.Args[1] == "report" || os.Args[1] == "review" {
		if len(os.Args) == 2 {
			fmt.Println("Correct usage for reporting or reviewing is, <binary> <report/review> <your issue/review>")
			return
		}
		payload["Payload_type"] = os.Args[1]
		payload["Username"] = username
		payload["Body"] = os.Args[2:]
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err = enc.Encode(&payload)
		conn.Write(buf.Bytes())
		fmt.Println("your request has been successfully submitted")
		return
	}
	downloadmanager.Run()
}
