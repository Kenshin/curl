package curl

import (
	"fmt"
	"os"
	"testing"
)

func TestCurl(t *testing.T) {

	/*
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
		code := New("http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir()+"/"+"node.exe")
		fmt.Printf("curl.New return code is %v\n", code)
	*/

	// multi download
	dl := Download{}
	dl.AddTask(Task{"http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir() + "/" + "node.exe", 0})
	dl.AddTask(Task{"http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node3.exe", os.TempDir() + "/" + "node3.exe", 0})
	dl.AddTask(Task{"https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "ChromeSetup.zip", os.TempDir() + "/" + "ChromeSetup.zip", 0})
	//dl.AddTask(Task{"http://npm.taobao.org/mirrors/node/v5.7.1/win-x64/node.exe", "node2.exe", os.TempDir() + "/" + "node2.exe", 0})
	//dl.AddTask(Task{"http://golangtc.com/static/go/1.6/go1.6.windows-amd64.zip", "windows-amd64.zip", os.TempDir() + "/" + "windows-amd64.zip", 0})

	/*
		dl[0] = Detail{"http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir() + "/" + "node.exe"}
		dl[1] = Detail{"http://npm.taobao.org/mirrors/node/v5.7.1/win-x64/node.exe", "node2.exe", os.TempDir() + "/" + "node2.exe"}
		dl[2] = Detail{"http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node3.exe", os.TempDir() + "/" + "node3.exe"}
		dl[3] = Detail{"http://npm.taobao.org/mirrors/node/v4.0.0/win-x86/node.exe", "node5.exe", os.TempDir() + "/" + "node5.exe"}
		dl[4] = Detail{"http://npm.taobao.org/mirrors/node/v4.1.0/win-x86/node.exe", "node6.exe", os.TempDir() + "/" + "node6.exe"}
		dl[5] = Detail{"http://npm.taobao.org/mirrors/node/v4.2.0/win-x64/node.exe", "node7.exe", os.TempDir() + "/" + "node7.exe"}
		dl[6] = Detail{"http://npm.taobao.org/mirrors/node/v4.3.0/win-x86/node.exe", "node8.exe", os.TempDir() + "/" + "node8.exe"}
		dl[7] = Detail{"http://npm.taobao.org/mirrors/node/v4.3.1/win-x64/node.exe", "node9.exe", os.TempDir() + "/" + "node9.exe"}
		dl[8] = Detail{"http://npm.taobao.org/mirrors/node/v4.3.2/win-x86/node.exe", "node10.exe", os.TempDir() + "/" + "node10.exe"}
		dl[9] = Detail{"http://npm.taobao.org/mirrors/node/v4.2.2/win-x64/node.exe", "node11.exe", os.TempDir() + "/" + "node11.exe"}
		dl[10] = Detail{"http://npm.taobao.org/mirrors/node/v4.2.1/win-x86/node.exe", "node12.exe", os.TempDir() + "/" + "node12.exe"}
	*/
	newDL, err := New(dl)
	fmt.Printf("curl.New return ld  is %v\n", newDL)
	fmt.Printf("curl.New return err is %v\n", err)

}
