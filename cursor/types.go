package cursor

const (
	getCursorPosition     string = "\x1b[6n"
	cursorPositionFormat  string = "\x1bs[%d;%dR" // x (line); y (column)
	setCursorPosition     string = "\x1b[%d;%dH"
	scrollUp              string = "\x1b[%dS"
	scrollDown            string = "\x1b[%dT"
	saveCursorPosition    string = "\x1b[s"
	restoreCursorPosition string = "\x1b[u"
	hideCursor            string = "\x1b[?25l"
	showCursor            string = "\x1b[?25h"
)

type CursorMovement interface {
	getString() string
}

type cursorMovement string

func (c cursorMovement) getString() string {
	return string(c)
}

const (
	CursorUp    cursorMovement = "\x1b[%dA"
	CursorDown  cursorMovement = "\x1b[%dB"
	CursorLeft  cursorMovement = "\x1b[%dC"
	CursorRight cursorMovement = "\x1b[%dD"
	ToNextLine  cursorMovement = "\x1b[%dE"
	ToPrevLine  cursorMovement = "\x1b[%dF"
)
