package curl

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestCurl(t *testing.T) {
	testGet()
	testReadLine()
	testSimpleNew()
	testMultiNewMode1()
	testMultiNewMode2()
}

func testGet() {
	code, res, _ := Get("http://npm.taobao.org/mirrors/node/latest/SHASUMS256.txt")
	if code != 0 {
		return
	}
	defer res.Body.Close()
}

func testReadLine() {
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
}

func testSimpleNew() {
	newDL, err := New("http://npm.taobao.org/mirrors/node/v0.10.26/node.exe")
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)
}

func testMultiNewMode1() {
	ts := Task{}
	ts1 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/visualhunt.json")
	ts2 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/02073.jpg")
	ts3 := ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/0207.jpg")
	newDL, err := New(ts1, ts2, ts3)
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)
}

func testMultiNewMode2() {
	ts := Task{}
	dl := Download{
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/visualhunt.json"),
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/02073.jpg", "1111.jpg"),
		ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/0207.jpg", "2222.jpg"),
	}
	/*
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
	*/
	dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "latest", "nodeeeeeeeeeeeeeeeeeeeeeeee.exe", os.TempDir()))
	dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "5.7.0", "node4.exe"))
	dl.AddTask(ts.New("https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "chrome 49.01", "ChromeSetup.zip"))
	newDL, err := New(dl)
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)
}
