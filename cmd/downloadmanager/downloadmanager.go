package downloadmanager

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	//UnixPathSeparator is wisott
	UnixPathSeparator = '/'
	//WindowsPathSeparator is wisott
	WindowsPathSeparator = '\\'
	//Unix is wisott
	Unix = "unix"
	//Windows is wisott
	Windows = "windows"
)

//Download struct is the driver of this program
type Download struct {
	TotalThreads int
	URL          string
	DownloadPath string
	OS           string
}

//determineOSCompileTime is wisott
func determineOSCompileTime() (string, error) {
	switch os.PathSeparator {
	case UnixPathSeparator:
		return Unix, nil
	case WindowsPathSeparator:
		return Windows, nil
	}
	return "", errors.New("Unable to identify the os")
}

// Get a new http request
func (d Download) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		d.URL,
		nil,
	)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "Vishal Download Manager")
	return r, nil
}

//DownloadFile actually downloads the file
func (d Download) DownloadFile() error {
	r, err := d.getNewRequest("HEAD")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Can't process, response is %v", resp.StatusCode)
	}
	fmt.Println("calculating size of the file")
	size, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		fmt.Println("unable to calculate size of file")
	}
	fmt.Printf("Size of the file is %v bytes", size)
	var threads = make([][2]int, d.TotalThreads)
	threadDownloadSize := size / d.TotalThreads
	fmt.Printf("Each size is %v bytes\n", threadDownloadSize)
	for i := range threads {
		if i == 0 {
			//starting byte of first thread
			threads[i][0] = 0
		} else {
			//starting byte of other threads
			threads[i][0] = threads[i-1][1] + 1
		}
		if i < d.TotalThreads-1 {
			//ending byte of other threads
			threads[i][1] = threads[i][0] + threadDownloadSize
		} else {
			//ending byte of last thread
			threads[i][1] = size - 1
		}
	}
	fmt.Println("threads ", threads)
	var wg sync.WaitGroup
	for i, t := range threads {
		wg.Add(1)
		go func(i int, t [2]int) {
			defer wg.Done()
			err = d.DownloadThread(i, t)
			if err != nil {
				panic(err)
			}
		}(i, t)
	}
	wg.Wait()
	return d.mergeFiles(threads)
}

func (d Download) mergeFiles(threads [][2]int) error {
	f, err := os.OpenFile(d.DownloadPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	for i := range threads {
		tmpFileName := fmt.Sprintf("section-%v.tmp", i)
		b, err := ioutil.ReadFile(tmpFileName)
		if err != nil {
			return err
		}
		n, err := f.Write(b)
		if err != nil {
			return err
		}
		err = os.Remove(tmpFileName)
		if err != nil {
			return err
		}
		fmt.Printf("%v bytes merged\n", n)
	}
	return nil
}

//checkUrl checks for the correctness of the url matching a regex pattern
func checkURL(urlString string) error {
	// timeout := 1 * time.Second

	// _, err := url.ParseRequestURI(urlString)
	// if err != nil {
	// 	fmt.Println("parse uri error")
	// 	return err
	// }

	// u, err := url.Parse(urlString)
	// if err != nil || u.Scheme == "" || u.Host == "" {
	// 	fmt.Println("parse uri error 2")
	// 	return err
	// }

	// _, err := net.DialTimeout("tcp", urlString, timeout)
	// if err != nil {
	// 	fmt.Println("parse uri error 3")
	// 	return err
	// }
	return nil
}

//DownloadThread downloads a section of a file concurrently
func (d Download) DownloadThread(i int, t [2]int) error {
	fmt.Println("downloading a specific thread")
	r, err := d.getNewRequest("GET")
	if err != nil {
		return err
	}
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", t[0], t[1]))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("Can't process, response is %v", resp.StatusCode)
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("section-%v.tmp", i), b, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

//getDownloadValues fills the download struct with the options
func getDownloadValues(d *Download) error {
	lenArgs := len(os.Args)
	var err error
	switch lenArgs {
	case 1:
		fmt.Println("correct usage <binary> <url> <filename>")
		return errors.New("Incorrect binary usage")
	case 2:
		d.URL = os.Args[1]
		d.DownloadPath, err = os.Getwd()
		d.DownloadPath = d.DownloadPath + "/vdownloadedfile"
	case 3:
		d.URL = os.Args[1]
		d.DownloadPath = os.Args[2]
	}
	err = checkURL(d.URL)
	if err != nil {
		fmt.Println("the url does not seem to be a valid url ", err)
		return err
	}
	if d.OS == "unix" {
		if d.DownloadPath[0] != '/' {
			//relative file path
			//need to append an absolute file path
		}
	}
	return nil
}

//Run driver of download manager
func Run() {
	cpus := runtime.NumCPU()
	if cpus == 1 {
		runtime.GOMAXPROCS(1)
	} else {
		runtime.GOMAXPROCS(cpus - 1)
	}
	startTime := time.Now()
	d := Download{}
	d.TotalThreads = 10
	var err error
	d.OS, err = determineOSCompileTime()
	if err != nil {
		panic(err)
	}
	err = getDownloadValues(&d)
	if err != nil {
		//printed necessary information on console to use it,
		// no need to panic the program now, just exit
		return
	}
	fmt.Println(d)
	fmt.Printf("Download started after %v seconds\n", time.Now().Sub(startTime).Seconds())
	err = d.DownloadFile()
	if err != nil {
		fmt.Println("download file error ", err)
	}
}
