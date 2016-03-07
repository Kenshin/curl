/*
Curl is Simple http download and readline lib by Golang. Vesion 0.0.2
Website https://github.com/kenshin/curl
Copyright (c) 2014-2016 Kenshin Wang <kenshin@ksria.com>
*/
package curl

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

const ESC = "\033["

var (
	curLine  int           = -1
	mutex    *sync.RWMutex = new(sync.RWMutex)
	errStack []curlError   = make([]curlError, 0)
)

type curlError struct {
	name    string
	code    int
	message interface{}
}

func (err curlError) Error() {
	fmt.Printf("Name  : %v\n", err.name)
	fmt.Printf("Code  : %v\n", err.code)
	fmt.Printf("Error : %v", err.message)
}

type Detail struct {
	Url  string
	Name string
	Dst  string
}

type Download []Detail

// Read line use callback Process
// Line by line to obtain content and line num
type processFunc func(content string, line int) bool

func (dl Download) Add(da Detail) Download {
	return append(dl, da)
}

func (dl Download) GetValues(key string) []string {
	var arr []string
	for i := 0; i < len(dl); i++ {
		v := reflect.ValueOf(dl[i]).FieldByName(key)
		arr = append(arr, v.String())
	}
	return arr
}

// sync
var wg sync.WaitGroup

// Get url method
//
//  url e.g. http://nodejs.org/dist/v0.10.0/node.exe
//
// Return code
//   0: success
//  -1: status code != 200
//
// Return res, err
//
// For example:
//  code, res, _ := curl.Get("http://nodejs.org/dist/")
//  if code != 0 {
//      return
//  }
//  defer res.Body.Close()
func Get(url string) (code int, res *http.Response, err error) {

	// get res
	res, err = http.Get(url)

	// err
	if err != nil {
		panic(err)
	}

	// check state code
	if res.StatusCode != 200 {
		fmt.Printf("URL [%v] an [%v] error occurred, please check.\n", url, res.StatusCode)
		return -1, res, err
	}

	return 0, res, err

}

// Read line from io.ReadCloser
//
// For example:
//  versionFunc := func(content string, line int) bool {
//    // TO DO
//    return false
//  }
//
//  if err := curl.ReadLine(res.Body, versionFunc); err != nil && err != io.EOF {
//    //TO DO
//  }
func ReadLine(body io.ReadCloser, process processFunc) error {

	var content string
	var err error
	var line int = 1

	// set buff
	buff := bufio.NewReader(body)

	for {
		content, err = buff.ReadString('\n')

		if line > 1 && (err != nil || err == io.EOF) {
			break
		}

		if ok := process(content, line); ok {
			break
		}

		line++
	}

	return err
}

// Download method
//
// Parameter
//  url : download url e.g. http://nodejs.org/dist/v0.10.0/node.exe
//  name: download file name e.g. node.exe
//  dst : download path
//
// Return code
//   0: success
//  -2: create file error.
//  -3: download node.exe size error.
//  -4: content length = -1.
//  -5: panic error.
//  -6: curl.New() parameter type error.
//  -7: curl.New() parameter type error.
//
// For example:
//  curl.New("http://nodejs.org/dist/", "0.10.28", "v0.10.28")
//
//  Console show:
//
//  Start download [0.10.28] from http://nodejs.org/dist/.
//  node.exe: 70% [==============>__________________] 925ms
//  End download.
//
func New(args ...interface{}) int {
	var (
		code, count int = 0, 0
		dl          Download
	)

	if len(args) == 3 {
		count = 1
		dl = dl.Add(Detail{args[0].(string), args[1].(string), args[2].(string)})
	} else if len(args) == 1 {
		if v, ok := args[0].(Download); !ok {
			return -6
		} else {
			dl = v
			count = len(dl)
		}
	} else {
		return -6
	}

	fmt.Printf("Start download [%v].\n%v", strings.Join(dl.GetValues("Name"), ", "))

	wg.Add(count)
	for i := 0; i < count; i++ {
		progressbar(dl[i].Name, time.Now(), 0, "\n")
		go func(dl Download, num int) {
			code = download(dl[num].Url, dl[num].Name, dl[num].Dst, num, count)
			wg.Done()
		}(dl, i)
	}
	wg.Wait()

	curDown(count - curLine)
	for _, v := range errStack {
		fmt.Println("\n-------- Error Message detail --------")
		v.Error()
	}
	fmt.Println("\r\n--------\nEnd download.")

	return code
}

func download(url, name, dst string, line, max int) int {
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(curlError); ok {
				errStack = append(errStack, v)
			} else {
				errStack = append(errStack, curlError{name, -5, err})
			}
			curStack(line, max)
			empty := strings.Repeat(" ", 100)
			fmt.Printf("\r%v download error.%v", name, empty)
		}
	}()

	// get url
	code, res, err := Get(url)
	if code != 0 {
		return code
	}
	defer res.Body.Close()

	// create file
	file, createErr := os.Create(dst)
	if createErr != nil {
		panic(curlError{name, -2, "Create file error, Error: " + createErr.Error()})
		return -2
	}
	defer file.Close()

	// verify content length
	if res.ContentLength == -1 {
		panic(curlError{name, -4, "Download content length is -1."})
		return -4
	}

	start := time.Now()
	buf := make([]byte, res.ContentLength)
	var m float32
	for {
		n, err := res.Body.Read(buf)
		if n == 0 && err.Error() == "EOF" {
			break
		}
		if err != nil && err.Error() != "EOF" {
			panic(curlError{name, -7, "Download size error, Error: ." + err.Error()})
			return -7
		}
		m = m + float32(n)
		i := int(m / float32(res.ContentLength) * 50)
		file.WriteString(string(buf[:n]))

		func(name string, start time.Time, i, line, max int) {
			curStack(line, max)
			progressbar(name, start, i, "")
		}(name, start, i, line, max)
	}

	// valid download exe
	fi, err := file.Stat()
	if err == nil {
		if fi.Size() != res.ContentLength {
			panic(curlError{name, -3, "Downlaod size verify error, please check your network."})
			return -3
		}
	}
	return 0
}

/*
 name: 70% [==============>__________________] 925ms
*/
func progressbar(name string, start time.Time, i int, suffix string) {
	h := strings.Repeat("=", i) + ">" + strings.Repeat("_", 50-i)
	d := time.Now().Sub(start)
	fmt.Printf("\r"+name+": "+"%.0f%% [%s] %v"+suffix, float32(i)/50*100, h, time.Duration(d.Seconds())*time.Second)
}

func curStack(line, max int) {
	mutex.Lock()
	switch {
	case curLine == -1:
		curUp(max - line)
	case line < curLine:
		curUp(line - curLine)
	case line > curLine:
		curDown(curLine - line)
	}
	if curLine != line {
		curLine = line
	}
	mutex.Unlock()
}

func curUp(i int) {
	fmt.Printf(ESC+"%dA", i)
}

func curDown(i int) {
	fmt.Printf(ESC+"%dB", i)
}
