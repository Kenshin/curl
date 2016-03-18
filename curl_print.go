package curl

import (
	"fmt"
	"strings"
	"time"
)

func header(dl *Download) {
	fmt.Printf("Start download [%v].\n%v", strings.Join((*dl).GetValues("Title"), ", "))
}

func footer() {
	fmt.Println("\r--------\nEnd download.")
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
