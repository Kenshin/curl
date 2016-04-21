package curl

import (
	"os"
	"strings"
)

/*
 resolve arguments and create Download struct
*/
func parseArgs(args ...interface{}) (int, Download) {
	dl := Download{}
	if len(args) == 0 {
		panic(CurlError{"curl.New()", -6, "curl.New() parameter type error."})
	} else {
		switch args[0].(type) {
		case string:
			url, title, name, dst := safeArgs(args...)
			dl.AddTask(Task{url, title, name, dst, 0})
		case Task:
			for _, v := range args {
				dl.AddTask(v.(Task))
			}
		case Download:
			dl = args[0].(Download)
		}
	}
	return len(dl), dl
}

/*

 resolve arguments and create url, title, name, dst

 arg1: url
 arg2: title
 arg3: name
 arg4: dst

*/
func safeArgs(args ...interface{}) (url, title, name, dst string) {
	url = args[0].(string)
	switch len(args) {
	case 1:
		names := strings.Split(url, "/")
		title = names[len(names)-1:][0]
		name = title
		dst, _ = os.Getwd()
	case 2:
		title = args[1].(string)
		name = title
		dst, _ = os.Getwd()
	case 3:
		title, name = args[1].(string), args[2].(string)
		dst, _ = os.Getwd()
	case 4:
		title, name, dst = args[1].(string), args[2].(string), args[3].(string)
	}
	return
}

/*
 verify path is folder
*/
func isDirExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return true
	}
}

/*
 verify Content-Type is binary
*/
func isBodyBytes(content string) (isBytes bool) {
	if strings.Index(content, "json") != -1 {
		isBytes = false
	} else if strings.Index(content, "text") != -1 {
		isBytes = false
	} else if strings.Index(content, "application") != -1 {
		isBytes = true
	}
	return
}

/*
 get max len(title) num
*/
func maxTitleLength(titles []string) int {
	max := 0
	for _, v := range titles {
		if len(v) > max {
			max = len(v)
		}
	}
	if max > 15 {
		max = 15
	}
	return max
}

/*
 when title out range of max(15) format xxxx...
*/
func safeTitle(title string) string {
	h := ""
	if len(title) > 15 {
		title = title[:12] + "..."
	} else if len(title) <= maxNameLen {
		h = strings.Repeat(" ", maxNameLen-len(title))
	}
	return h + title + ":"
}

/*
 when dst the end of no '/', return format dst
*/
func safeDst(dst string) string {
	if !strings.HasSuffix(dst, "/") {
		dst += "/"
	}
	return dst
}
