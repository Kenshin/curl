// +build !windows

package curl

import "fmt"

func curUp(i int) {
	fmt.Printf("\r\033[%dA", i)
}

func curDown(i int) {
	fmt.Printf("\r\033[%dB", i)
}
