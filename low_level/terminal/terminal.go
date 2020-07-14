package terminal

// #include "terminal.h"
import "C"

func GetTerminalSize() (int, int) {
	size := C.getTerminalSize()
	return int(size.width), int(size.height)
}
