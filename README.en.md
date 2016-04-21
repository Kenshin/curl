curl - Multiple download lib with CLI by Golang  
================================
[![Build Status](https://api.travis-ci.org/Kenshin/curl.svg?branch=master)](https://travis-ci.org/Kenshin/curl)
[![Version][version-badge]][version-link]
[![Gowalker][gowalker-badge]][gowalker-link]
[![Godoc][godoc-badge]][godoc-link]
[![Gitter][gitter-badge]][gitter-link]
[![Slack][slack-badge]][slack-link]
[![Jianliao][jianliao-badge]][jianliao-link]  

`curl` `Multiple download` lib with CLI by `Golang`, `content-type` include: application( `exe`, `jpg`, etc. ), text( `txt`, `json`, etc. ).  

![Multi-download](http://i.imgur.com/BRb7vm1.gif)

Document
---
[English](https://github.com/kenshin/curl/blob/master/README.en.md) | [繁體中文](https://github.com/kenshin/curl/blob/master/README.tw.md)

Support
---
* Mac OS
* Linux
* Windows ( usage `kernel32.dll` and `SetConsoleCursorPosition` function )

Installation
---
`go get -u github.com/Kenshin/curl`

Usage
---
`import "github.com/Kenshin/curl"`

Getting Started
---
##### Read line
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

##### Simple download
```
newDL, err := curl.New("http://npm.taobao.org/mirrors/node/v0.10.26/node.exe")
fmt.Printf("curl.New return ld  is %v\n", newDL)
fmt.Printf("curl.New return err is %v\n", err)
```
![Simple-download](http://i.imgur.com/bNBJ2kG.png)

##### Multi download
```
// mode 1
ts := curl.Task{}
ts1 := ts.New("http://xxx.xxx.xxx/node/latest/node.exe", "node.exe")
ts2 := ts.New("http://xxx.xxx.xxx/node/v4.0.0/win-x64/node.exe", "node40.exe")
ts3 := ts.New("http://xxx.xxx.xxx/node/v4.1.0/win-x64/node.exe", "node41.exe")
ts4 := ts.New("http://xxx.xxx.xxx/node/v4.2.0/win-x64/node.exe", "node42.exe")
ts5 := ts.New("http://xxx.xxx.xxx/node/v4.3.0/win-x64/node41.exe", "node43.exe")
newDL, err := curl.New(ts1, ts2, ts3, ts4, ts5, ts6)

fmt.Printf("curl.New return ld  is %v\n", newDL)
fmt.Printf("curl.New return err is %v\n", err)

// mode 2
dl := curl.Download {
    ts.New("http://7x2xql.com1.z0.glb.clouddn.com/visualhunt.json"),
    ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/02073.jpg"),
    ts.New("http://7x2xql.com1.z0.glb.clouddn.com/holiday/0207.jpg"),
}
dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/latest/node.exe", "nodeeeeeeeeeeeeeeeeeeeeeeee.exe", os.TempDir()))
dl.AddTask(ts.New("http://npm.taobao.org/mirrors/node/v5.7.0/win-x64/node.exe", "node4.exe", os.TempDir()))
dl.AddTask(ts.New("https://www.google.com/intl/zh-CN/chrome/browser/?standalone=1&extra=devchannel&platform=win64", "ChromeSetup.zip", os.TempDir()))
newDL, err := curl.New(dl)

fmt.Printf("curl.New return ld  is %v\n", newDL)
fmt.Printf("curl.New return err is %v\n", err)
```
![Multi-download](http://i.imgur.com/BRb7vm1.gif)

##### Custom download process
![custom progress schematic](http://i.imgur.com/F5xjerv.jpg)
```
// npm like
curl.Options.Header = false
curl.Options.Footer = false
curl.Options.LeftEnd = ""
curl.Options.RightEnd = ""
curl.Options.Fill = "█"
curl.Options.Arrow = ""
curl.Options.Empty = "░"

newDL, err := New("http://npm.taobao.org/mirrors/node/v0.10.26/node.exe")

node.exe: 100% ███████████████████████████████████████░░░░░░░░ 4s
```
![custom download progressbar](http://i.imgur.com/qokcgfB.gif)

About
---
* [CHANGELOG](https://github.com/kenshin/curl/blob/master/CHANGELOG.md)
* [Contact](http://kenshin.wang/) | [Email](kenshin@ksria.com) | [Twitter](https://twitter.com/wanglei001)
* [Feedback](https://github.com/kenshin/curl/issues)

CHANGELOG
---
* **2016-03-10, Version `0.0.4` support:**
    * Add multi download.
    * Add custom progress.
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
[gowalker-badge]:   https://img.shields.io/badge/go_walker-document-green.svg
[gowalker-link]:    http://gowalker.org/github.com/kenshin/curl
[godoc-badge]:      https://godoc.org/github.com/kenshin/curl?status.svg
[godoc-link]:       https://godoc.org/github.com/kenshin/curl
[version-badge]:    https://img.shields.io/badge/lastest_version-0.0.4-blue.svg
[version-link]:     https://github.com/kenshin/curl/releases
[travis-badge]:     https://travis-ci.org/Kenshin/curl.svg?branch=master
[travis-link]:      https://travis-ci.org/Kenshin/curl
[gitter-badge]:     https://badges.gitter.im/kenshin/curl.svg
[gitter-link]:      https://gitter.im/kenshin/curl?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge
[slack-badge]:      https://img.shields.io/badge/chat-slack-orange.svg
[slack-link]:       https://curl.slack.com/
[jianliao-badge]:   https://img.shields.io/badge/chat-jianliao-yellowgreen.svg
[jianliao-link]:    https://guest.jianliao.com/rooms/76dce8b01v
[license-badge]:    https://img.shields.io/github/license/mashape/apistatus.svg
[license-link]:     https://opensource.org/licenses/MIT
