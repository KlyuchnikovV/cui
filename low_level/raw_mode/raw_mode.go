package raw_mode

// #include "./raw_mode.h"
import "C"

func EnableRawMode() {
	C.enableRawMode()
}

func DisableRawMode() {
	C.disableRawMode()
}
