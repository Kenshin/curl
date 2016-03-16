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
	newDL, err := New("http://npm.taobao.org/mirrors/node/v0.10.26/node.exe")
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)

	ts := Task{}
	ts1 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/visualhunt.json")
	ts2 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/02073.jpg")
	ts3 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/0207.jpg")
	newDL, err = New(ts1, ts2, ts3)
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)

	// multi download
	dl := Download{
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/visualhunt.json"),
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/02073.jpg"),
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/0207.jpg"),
	}
	dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "nodeeeeeeeeeeeeeeeeeeeeeeee.exe", os.TempDir()))
	dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node4.exe", os.TempDir()))
	dl.AddTask(ts.New("https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "ChromeSetup.zip", os.TempDir()))
	newDL, err = New(dl)

	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)

	/*
		dl := Download{}

		url := "http://npm.taobao.org/mirrors/node/latest/node.exe"
		parseArgs(&dl, errStack, url)

		ts := Task{}
		ts1 := ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node4.exe", os.TempDir())
		ts2 := ts.New("http://npm.taobao.org/mirrors/node/v5.7.2/win-x64/node.exe", "node2.exe", os.TempDir())
		ts3 := ts.New("http://npm.taobao.org/mirrors/node/v5.7.3/win-x64/node.exe", "node3.exe", os.TempDir())
		parseArgs(&dl, errStack, ts1, ts2, ts3)

		ts := new(Task)
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir()))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node4.exe", os.TempDir()))
		dl.AddTask(ts.New("https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "ChromeSetup.zip", os.TempDir()))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "nodeeeeeeeeeeeeeeeeeeeeeeee.exe", os.TempDir()))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node4.exe", os.TempDir()))
		dl.AddTask(ts.New("https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "ChromeSetup.zip", os.TempDir()))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x86/node.exe", "node40.exe"))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v4.1.1/win-x86/node.exe", "node41.exe"))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v4.4.0/win-x86/node.exe", "node42.exe"))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v4.3.0/win-x86/node.exe", "node43.exe"))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v4.2.0/win-x86/node.exe", "node44.exe"))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v4.1.0/win-x86/node.exe", "node45.exe"))
		dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v4.0.0/win-x86/node.exe", "node46.exe"))
		dl.AddTask(ts.New("http://7x2xql.com1.z0.glb.clouddn.com/visualhunt.json"))
		dl.AddTask(ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/02073.jpg"))
		dl.AddTask(ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/0207.jpg"))
		parseArgs(&dl, dl)
	*/

}
