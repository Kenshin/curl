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

// Read line use callback Process
// Line by line to obtain content and line num
type processFunc func(content string, line int) bool

type detail struct {
	url  string
	name string
	dst  string
}
type Download map[int]detail

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
//  -2: create file error
//  -3: download node.exe error
//  -4: content length = -1
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
		code, count int      = 0, 0
		dl          Download = make(Download)
	)

	if len(args) == 3 {
		count = 1
		dl[count-1] = detail{args[0].(string), args[1].(string), args[2].(string)}
	}

	wg.Add(count)

	fmt.Printf("Start download [%v].\n%v", strings.Join(dl.GetValues("name"), ","))

	for i := 0; i < count; i++ {
		go func(dl Download, num int) {
			urls := dl[num]
			code = download(urls.url, urls.name, urls.dst)
			wg.Done()
		}(dl, i)
	}

	wg.Wait()
	fmt.Println("\nEnd download.")

	return code
}

func download(url, name, dst string) int {
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Sprintf("CURL Error: Download %v from %v an error has occurred. \nError: %v", name, url, err)
			panic(msg)
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
		fmt.Println("Create file error, Error: " + createErr.Error())
		return -2
	}
	defer file.Close()

	// verify content length
	if res.ContentLength == -1 {
		fmt.Printf("Download %v fail from %v.\n", name, url)
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
			panic(err)
		}
		m = m + float32(n)
		i := int(m / float32(res.ContentLength) * 50)
		file.WriteString(string(buf[:n]))
		progressbar(name, start, i)
	}

	// valid download exe
	fi, err := file.Stat()
	if err == nil {
		if fi.Size() != res.ContentLength {
			fmt.Printf("Error: Downlaod [%v] size error, please check your network.\n", name)
			return -3
		}
	}
	return 0
}

/*
 name: 70% [==============>__________________] 925ms
*/
func progressbar(name string, start time.Time, i int) {
	h := strings.Repeat("=", i) + ">" + strings.Repeat("_", 50-i)
	d := time.Now().Sub(start)
	fmt.Printf("\r"+name+": "+"%.0f%% [%s] %v", float32(i)/50*100, h, time.Duration(d.Seconds())*time.Second)
}
