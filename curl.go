/*
Curl is Simple http download and readline lib by Golang. Vesion 0.0.4
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

const esc = "\033["

var (
	wg      sync.WaitGroup
	curLine int           = -1
	mutex   *sync.RWMutex = new(sync.RWMutex)
)

// Curl Error struct
type CurlError struct {
	name    string      // Task struct Name
	code    int         // Task struct Code
	message interface{} // Error message
}

// Print Error
func (err CurlError) Error() {
	fmt.Printf("Name  : %v\n", err.name)
	fmt.Printf("Code  : %v\n", err.code)
	fmt.Printf("Error : %v", err.message)
}

// Task struct
type Task struct {
	Url  string
	Name string
	Dst  string
	Code int
}

// Receive url, name and dst
// Retruns New Task
func (ts Task) New(url, name, dst string) Task {
	ts.Url, ts.Name, ts.Dst = url, name, dst
	return ts
}

type Download struct {
	tasks []Task
}

// Append Download task arrray
func (dl *Download) AddTask(ts Task) {
	dl.tasks = append(dl.tasks, ts)
}

// Get Download struct values by key
func (dl Download) GetValues(key string) []string {
	var arr []string
	for i := 0; i < len(dl.tasks); i++ {
		v := reflect.ValueOf(dl.tasks[i]).FieldByName(key)
		arr = append(arr, v.String())
	}
	return arr
}

// Read line use callback Process
// Line by line to obtain content and line num
type processFunc func(content string, line int) bool

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

/*
   Download method

   Parameter:
    simple download model:
        url : download url e.g. http://nodejs.org/dist/v0.10.0/node.exe
        name: download file name e.g. node.exe
        dst : download path
    multi download model:
    Download struct

   Return code:
     0: success
    -2: create file error.
    -3: download node.exe size error.
    -4: content length = -1.
    -5: panic error.
    -6: curl.New() parameter type error.
    -7: Download size error.

   Return:
    dl( Download struct )
    err( []CurlError array)

   For example:

    // simple download
    dl, err := curl.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir()+"/"+"node.exe")

    // multi download
    dl := curl.Download{}
    ts := new(curl.Task)
    dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir()+"/"+"node.exe"))
    dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node3.exe", os.TempDir()+"/"+"node3.exe"))
    dl.AddTask(ts.New("https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "ChromeSetup.zip", os.TempDir()+"/"+"ChromeSetup.zip"))
    newDL, err := New(dl)

   Console show:
    Start download [aaa, bbb, node, npm].
         aaa: 70% [==============>__________________] 925ms
         bbb: 10% [===>_____________________________] 2s
        node: 100% [===============================>] 10s
         npm: download error.
    End download.
*/
func New(args ...interface{}) (dl Download, errStack []CurlError) {
	var count int = 0
	curLine = -1
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(CurlError); ok {
				errStack = append(errStack, v)
			} else {
				errStack = append(errStack, CurlError{"curl.New()", -5, err})
			}
		}
	}()

	if len(args) == 3 {
		count = 1
		dl.AddTask(Task{args[0].(string), args[1].(string), args[2].(string), 0})
	} else if len(args) == 1 {
		if v, ok := args[0].(Download); !ok {
			panic(CurlError{"curl.New()", -6, "curl.New() parameter type error."})
		} else {
			dl = v
			count = len(dl.tasks)
		}
	} else {
		panic(CurlError{"curl.New()", -6, "curl.New() parameter type error."})
	}

	fmt.Printf("Start download [%v].\n%v", strings.Join(dl.GetValues("Name"), ", "))

	wg.Add(count)
	for i := 0; i < count; i++ {
		progressbar(dl.tasks[i].Name, time.Now(), 0, "\n")
		go func(dl Download, num int) {
			download(&dl.tasks[num], num, count, &errStack)
			wg.Done()
		}(dl, i)
	}
	wg.Wait()

	curDown(count - curLine)
	fmt.Println("\r--------\nEnd download.")

	return
}

func download(ts *Task, line, max int, errStack *[]CurlError) {
	url, name, dst := ts.Url, ts.Name, ts.Dst
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(CurlError); ok {
				*errStack = append(*errStack, v)
				ts.Code = v.code
			} else {
				*errStack = append(*errStack, CurlError{name, -5, err})
				ts.Code = -5
			}
			curStack(line, max)
			empty := strings.Repeat(" ", 100)
			fmt.Printf("\r%v download error.%v", name, empty)
		}
	}()

	// get url
	code, res, err := Get(url)
	if code == -1 {
		panic(CurlError{name, -1, "curl.Get() error, Error: " + err.Error()})
	}
	defer res.Body.Close()

	// create file
	file, createErr := os.Create(dst)
	if createErr != nil {
		panic(CurlError{name, -2, "Create file error, Error: " + createErr.Error()})
	}
	defer file.Close()

	// verify content length
	if res.ContentLength == -1 {
		panic(CurlError{name, -4, "Download content length is -1."})
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
			panic(CurlError{name, -7, "Download size error, Error: ." + err.Error()})
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
			panic(CurlError{name, -3, "Downlaod size verify error, please check your network."})
		}
	}
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
	fmt.Printf(esc+"%dA", i)
}

func curDown(i int) {
	fmt.Printf(esc+"%dB", i)
}
