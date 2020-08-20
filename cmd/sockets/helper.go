package sockets

import (
	"errors"
	"net/http"
	"os"
	"runtime"
	"strings"
)

//OS is the operating system of the target system
type OS string

type User struct {
	//Uid user id
	Uid string
	// Gid is the primary group ID.
	// On POSIX systems, this is a decimal number representing the gid.
	// On Windows, this is a SID in a string format.
	// On Plan 9, this is the contents of /dev/user.
	Gid string
	// Username is the login name.
	Username string
	// Name is the user's real or display name.
	// It might be blank.
	// On POSIX systems, this is the first (or only) entry in the GECOS field
	// list.
	// On Windows, this is the user's display name.
	// On Plan 9, this is the contents of /dev/user.
	Name    string
	HomeDir string
}

//Target contains information about the system on which binary is deployed
type Target struct {
	OS   OS
	Name string
	User User
}

// Get a new http request
func getNewRequest(method, url string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		url,
		nil,
	)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Vishal Download Manager")
	return r, nil
}

const (
	//UnixPathSeparator is wisott
	UnixPathSeparator = '/'
	//WindowsPathSeparator is wisott
	WindowsPathSeparator = '\\'
	//Unix is wisott
	unix = "darwin"
	//Windows is wisott
	windows = "windows"
	//Linux is wisott
	linux = "linux"
)

const (
	//Unix OS
	Unix OS = unix
	//Windows OS
	Windows OS = windows
	//Linux OS
	Linux OS = linux
)

//determineOSCompileTime is wisott
func determineOSCompileTime() (OS, error) {
	switch os.PathSeparator {
	case UnixPathSeparator:
		return Unix, nil
	case WindowsPathSeparator:
		return Windows, nil
	}
	return "", errors.New("Unable to identify the os")
}

//determineOSCompileTime is wisott
func determineOSRunTime() (OS, error) {
	switch strings.ToLower(runtime.GOOS) {
	case unix:
		return Unix, nil
	case linux:
		return Linux, nil
	case windows:
		return Windows, nil
	}
	return "", errors.New("Unable to detect runtime os")
}

func fetchWindowsDetails() {

}

func fetchLinuxDetails() {

}

func fetchUnixDetails() {

}
