package curl

import (
	"fmt"
	"syscall"
	"unsafe"
)

type (
	SHORT int16
	WORD  uint16

	SMALL_RECT struct {
		Left   SHORT
		Top    SHORT
		Right  SHORT
		Bottom SHORT
	}

	COORD struct {
		X SHORT
		Y SHORT
	}

	CONSOLE_SCREEN_BUFFER_INFO struct {
		Size              COORD
		CursorPosition    COORD
		Attributes        WORD
		Window            SMALL_RECT
		MaximumWindowSize COORD
	}
)

var (
	kernel32DLL                    = syscall.NewLazyDLL("kernel32.dll")
	getConsoleScreenBufferInfoProc = kernel32DLL.NewProc("GetConsoleScreenBufferInfo")
	setConsoleCursorPositionProc   = kernel32DLL.NewProc("SetConsoleCursorPosition")
	stdoutHandle                   = getStdHandle(syscall.STD_OUTPUT_HANDLE)
	position                       = COORD{1, -1}
)

func getCursorPosition() {
	if position.Y == -1 {
		var info, err = getConsoleScreenBufferInfo(stdoutHandle)
		if err != nil {
			panic("could not get console screen buffer info")
		}
		position = info.CursorPosition
	}
}

func getError(r1, r2 uintptr, lastErr error) error {
	if r1 == 0 {
		if lastErr != nil {
			return lastErr
		}
		return syscall.EINVAL
	}
	return nil
}

func getStdHandle(stdhandle int) uintptr {
	handle, err := syscall.GetStdHandle(stdhandle)
	if err != nil {
		panic(fmt.Errorf("could not get standard io handle %d", stdhandle))
	}
	return uintptr(handle)
}

func getConsoleScreenBufferInfo(handle uintptr) (*CONSOLE_SCREEN_BUFFER_INFO, error) {
	var info CONSOLE_SCREEN_BUFFER_INFO
	if err := getError(getConsoleScreenBufferInfoProc.Call(handle, uintptr(unsafe.Pointer(&info)), 0)); err != nil {
		return nil, err
	}
	return &info, nil
}

func setConsoleCursorPosition(handle uintptr, position COORD) {
	if err := getError(setConsoleCursorPositionProc.Call(handle, uintptr(*(*int32)(unsafe.Pointer(&position))))); err != nil {
		panic(err)
	}
}

func curReset(i int) {
	var info, err = getConsoleScreenBufferInfo(stdoutHandle)
	if err != nil {
		panic("could not get console screen buffer info")
	}
	position = info.CursorPosition
	position.Y -= SHORT(i)
	setConsoleCursorPosition(stdoutHandle, position)
}

func curUp(i int) {
	getCursorPosition()
	position.Y -= SHORT(i)
	setConsoleCursorPosition(stdoutHandle, position)
}

func curDown(i int) {
	getCursorPosition()
	position.Y += SHORT(i)
	setConsoleCursorPosition(stdoutHandle, position)
}
