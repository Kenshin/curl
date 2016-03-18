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
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	wg         sync.WaitGroup
	curLine    int = -1
	maxNameLen int
	mutex      *sync.RWMutex = new(sync.RWMutex)
	count      int           = 0
)

type (
	// Curl Error struct
	CurlError struct {
		name    string      // Task struct Name
		code    int         // Task struct Code
		message interface{} // Error message
	}

	// Task struct
	Task struct {
		Url   string
		Title string
		Name  string
		Dst   string
		Code  int
	}

	// Task array
	Download []Task

	// Read line use callback Process
	// Line by line to obtain content and line num
	processFunc func(content string, line int) bool
)

// Print Error
func (err CurlError) Error() string {
	name := fmt.Sprintf("Name  : %v\n", err.name)
	code := fmt.Sprintf("Code  : %v\n", err.code)
	msg := fmt.Sprintf("Error : %v", err.message)
	return "\n" + name + code + msg
}

// Receive url, name and dst
// Retruns New Task
func (ts Task) New(args ...interface{}) Task {
	if len(args) == 0 {
		panic(CurlError{"curl.New()", -6, "curl.New() parameter type error."})
	} else {
		ts.Url, ts.Title, ts.Name, ts.Dst = safeArgs(args...)
	}
	return ts
}

// Append Download task arrray
func (dl *Download) AddTask(ts Task) {
	*dl = append(*dl, ts)
}

// Get Download struct values by key
func (dl Download) GetValues(key string) []string {
	var arr []string
	for i := 0; i < len(dl); i++ {
		v := reflect.ValueOf(dl[i]).FieldByName(key)
		arr = append(arr, v.String())
	}
	return arr
}

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
		return -5, res, CurlError{url, -5, err.Error()}
	}

	// check state code
	if res.StatusCode != 200 {
		s := fmt.Sprintf("%v an [%v] error occurred.", url, res.StatusCode)
		return -1, res, CurlError{url, -1, s}
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
		url  : download url e.g. http://nodejs.org/dist/v0.10.0/node.exe
		title: download task label.
		name : download file name e.g. node.exe
		dst  : download path
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
	-8: Write Content-Type:text error.

   Return:
	dl( []Task Download struct )
	err( []CurlError array)

   Console show:
	Start download [aaa, bbb, node, npm, ccccccc].
		aaa: 70% [==============>__________________] 925ms
		bbb: 10% [===>_____________________________] 2s
	   node: 100% [===============================>] 10s
		npm: download error.
	cccc...: 30% [=========>________________________] 2s
	End download.

   For example:

	// simple download
	newDL, err := New("http://npm.taobao.org/mirrors/node/v0.10.26/node.exe")
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)

	// multi download
	ts := Task{}
	ts1 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/visualhunt.json")
	ts2 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/02073.jpg")
	ts3 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/0207.jpg")
	newDL, err = New(ts1, ts2, ts3)
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)

	dl := Download{
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/visualhunt.json"),
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/02073.jpg"),
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/0207.jpg"),
	}
	dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "nodeeeeeeeeeeeeeeeeeeeeeeee.exe"))
	dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node4.exe"))
	dl.AddTask(ts.New("https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "Chrome 49"))
	newDL, err = New(dl)

	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)
*/
func New(args ...interface{}) (dl Download, errStack []CurlError) {
	curLine = -1
	count, dl = parseArgs(args...)
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(CurlError); ok {
				errStack = append(errStack, v)
			} else {
				errStack = append(errStack, CurlError{"curl.New()", -5, err})
			}
		}
	}()

	fmt.Printf("Start download [%v].\n%v", strings.Join(dl.GetValues("Title"), ", "))
	maxNameLen = maxTitleLength(dl.GetValues("Title"))

	wg.Add(count)
	for i := 0; i < count; i++ {
		progressbar(dl[i].Title, time.Now(), 0, "\n")
		go func(dl Download, num int) {
			download(&dl[num], num, count, &errStack)
			wg.Done()
		}(dl, i)
	}
	wg.Wait()

	curDown(count - curLine)
	fmt.Println("\r--------\nEnd download.")

	return
}

/*
 download ( text/binary ) and save it.
*/
func download(ts *Task, line, max int, errStack *[]CurlError) {
	url, title, name, dst := ts.Url, ts.Title, ts.Name, safeDst(ts.Dst)
	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(CurlError); ok {
				*errStack = append(*errStack, v)
				ts.Code = v.code
			} else {
				*errStack = append(*errStack, CurlError{url, -5, err})
				ts.Code = -5
			}
			curMove(line, max)
			msg := fmt.Sprintf("%v download error.", safeTitle(title))
			empty := strings.Repeat(" ", 80-len(msg))
			fmt.Printf("\r%v%v", msg, empty)
		}
	}()

	// get url
	code, res, err := Get(url)
	if code == -1 {
		panic(err)
	}
	defer res.Body.Close()

	// create dst
	if !isDirExist(dst) {
		if err := os.Mkdir(dst, 0777); err != nil {
			panic(CurlError{url, -2, "Create folder error, Error: " + err.Error()})
		}
	}

	// create file
	file, createErr := os.Create(dst + name)
	if createErr != nil {
		panic(CurlError{url, -2, "Create file error, Error: " + createErr.Error()})
	}
	defer file.Close()

	// verify content length
	if res.ContentLength == -1 && isBodyBytes(res.Header.Get("Content-Type")) {
		panic(CurlError{url, -4, "Download content length is -1."})
	}

	start := time.Now()
	if isBodyBytes(res.Header.Get("Content-Type")) {
		buf := make([]byte, res.ContentLength)
		var m float32
		for {
			n, err := res.Body.Read(buf)
			if n == 0 && err.Error() == "EOF" {
				break
			}
			if err != nil && err.Error() != "EOF" {
				panic(CurlError{url, -7, "Download size error, Error: ." + err.Error()})
			}
			m = m + float32(n)
			i := int(m / float32(res.ContentLength) * 50)
			file.WriteString(string(buf[:n]))

			func(title string, start time.Time, i, line, max int) {
				curMove(line, max)
				progressbar(title, start, i, "")
			}(title, start, i, line, max)
		}

		// valid download exe
		fi, err := file.Stat()
		if err == nil {
			if fi.Size() != res.ContentLength {
				panic(CurlError{url, -3, "Downlaod size verify error, please check your network."})
			}
		}
	} else {
		if bytes, err := ioutil.ReadAll(bufio.NewReader(res.Body)); err != nil {
			panic(CurlError{url, -8, err.Error()})
		} else {
			file.Write(bytes)
			curMove(line, max)
			progressbar(title, start, 50, "")
		}
	}
}

/*
 title: 70% [==============>__________________] 925ms
*/
func progressbar(title string, start time.Time, i int, suffix string) {
	h := strings.Repeat("=", i) + ">" + strings.Repeat("_", 50-i)
	d := time.Now().Sub(start)
	s := fmt.Sprintf("%v %.0f%% [%s] %v", safeTitle(title), float32(i)/50*100, h, time.Duration(d.Seconds())*time.Second)
	if len(s) > 80 {
		s = s[:80]
	}
	e := strings.Repeat(" ", 80-len(s))
	fmt.Printf("\r%v%v%v", s, e, suffix)
}

/*
  cursor move( up and down )
*/
func curMove(line, max int) {
	mutex.Lock()
	switch {
	case curLine == -1:
		curReset(max - line)
	case line < curLine:
		curUp(curLine - line)
	case line > curLine:
		curDown(line - curLine)
	}
	if curLine != line {
		curLine = line
	}
	mutex.Unlock()
}
