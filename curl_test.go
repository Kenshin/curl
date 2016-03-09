package curl

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestCurl(t *testing.T) {

	// curl.Get
	code, res, _ := Get("http://npm.taobao.org/mirrors/node/latest/SHASUMS256.txt")
	if code != 0 {
		return
	}
	// close
	defer res.Body.Close()

	// parse callback
	processFunc := func(content string, line int) bool {
		fmt.Printf("line is %v, content is %v", line, content)
		return false
	}
	// relad line
	if err := ReadLine(res.Body, processFunc); err != nil && err != io.EOF {
		fmt.Println(err)
	}

	// simple download
	newDL, err := New("http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir())
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)

	// multi download
	dl := Download{}
	ts := new(Task)
	dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir()))
	dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node4.exe", os.TempDir()))
	dl.AddTask(ts.New("https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "ChromeSetup.zip", os.TempDir()))

	newDL, err = New(dl)
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)

}
