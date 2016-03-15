Curl: Multiple download lib with CLI by Golang [![Build Status](https://api.travis-ci.org/Kenshin/curl.svg?branch=master)](https://travis-ci.org/Kenshin/curl)
================================

Documentation
---
[![Gowalker][gowalker-badge]][gowalker-link]
[![Godoc][godoc-badge]][godoc-link]

Installation
---
`go get -u github.com/Kenshin/curl`

Usage
---
`import "github.com/Kenshin/curl"`

Example
---
##### Read line:
```
// curl.Get
code, res, _ := curl.Get("http://npm.taobao.org/mirrors/node/latest/SHASUMS256.txt")
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

// read line
if err := ReadLine(res.Body, processFunc); err != nil && err != io.EOF {
    fmt.Println(err)
}
```
![ReadLine](http://i.imgur.com/7kUdIpE.png)

##### Simple download:
```
newDL, err := curl.New("http://npm.taobao.org/mirrors/node/v0.10.26/node.exe")
fmt.Printf("curl.New return ld  is %v\n", newDL)
fmt.Printf("curl.New return err is %v\n", err)
```
![Simple-download](http://i.imgur.com/bNBJ2kG.png)

##### Multi download:
```
// mode 1
ts := new(curl.Task)
ts1 := ts.New("http://xxx.xxx.xxx/node/latest/node.exe", "node.exe")
ts2 := ts.New("http://xxx.xxx.xxx/node/v4.0.0/win-x64/node.exe", "node40.exe")
ts3 := ts.New("http://xxx.xxx.xxx/node/v4.1.0/win-x64/node.exe", "node41.exe")
ts4 := ts.New("http://xxx.xxx.xxx/node/v4.2.0/win-x64/node.exe", "node42.exe")
ts5 := ts.New("http://xxx.xxx.xxx/node/v4.3.0/win-x64/node41.exe", "node43.exe")
newDL, err := New(ts1, ts2, ts3, ts4, ts5, ts6)

// mode 2
dl := curl.Download{}
ts := new(curl.Task)
dl.AddTask(ts.New("http://xxx.xxx.xxx/node/latest/node.exe", "node.exe"))
dl.AddTask(ts.New("http://xxx.xxx.xxx/node/v4.0.0/win-x64/node.exe", "node40.exe"))
dl.AddTask(ts.New("http://xxx.xxx.xxx/node/v4.1.0/win-x64/node.exe", "node41.exe"))
dl.AddTask(ts.New("http://xxx.xxx.xxx/node/v4.2.0/win-x64/node.exe", "node42.exe"))
dl.AddTask(ts.New("http://xxx.xxx.xxx/node/v4.3.0/win-x64/node41.exe", "node43.exe"))
newDL, err := New(dl)

fmt.Printf("curl.New return ld  is %v\n", newDL)
fmt.Printf("curl.New return err is %v\n", err)
```
![Multi-download](http://i.imgur.com/BRb7vm1.gif)

Support
---
* Mac OS
* Linux
* Windows ( usage `kernel32.dll` and `SetConsoleCursorPosition` fuction )

Help
---
* Email: <kenshin@ksria.com>
* [Github issue](https://github.com/Kenshin/curl/issues/new)

CHANGELOG
---
* **2016-03-10, Version `0.0.4` support:**
    * Add multi download.
    * Rework `curl.New` function.
    * Adapter Go 1.6.

* **2016-03-05, Version `0.0.3` support:**
    * Add beautiful dowload print.

* **2014-07-10, Version `0.0.2` support:**
    * Adapter Go 1.3.

* **2014-05-28, Version `0.0.1` support:**
    * New
    * Get
    * Readline

LICENSE
---
[![license-badge]][license-link]

<!-- Link -->
[gowalker-badge]:   https://img.shields.io/badge/go_walker-documentation-green.svg
[gowalker-link]:    http://gowalker.org/github.com/kenshin/curl
[godoc-badge]:      https://godoc.org/github.com/kenshin/curl?status.svg
[godoc-link]:       https://godoc.org/github.com/kenshin/curl
[license-badge]:    https://img.shields.io/github/license/mashape/apistatus.svg
[license-link]:     https://opensource.org/licenses/MIT
