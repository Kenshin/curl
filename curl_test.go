package curl

import (
	"fmt"
	//"io"
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

	// simple download
	code := New("http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir()+"/"+"node.exe")
	fmt.Printf("curl.New return code is %v\n", code)
	/*
		// multi download
		dl := make(Download)
		dl[0] = detail{"http://npm.taobao.org/mirrors/node/latest/node.exe", "node.exe", os.TempDir() + "/" + "node.exe"}
		dl[1] = detail{"http://npm.taobao.org/mirrors/node/v5.7.1/win-x64/node.exe", "node2.exe", os.TempDir() + "/" + "node2.exe"}
		code := New(dl)
		fmt.Printf("curl.New return code is %v\n", code)
	*/

}
