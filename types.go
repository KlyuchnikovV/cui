package cui

type renderRequest []byte

func (r renderRequest) GetMessage() interface{} {
	return r
}

const (
	getCursorPosition    string = "\x1b[6n"
	cursorPositionFormat string = "\\x1b\\[([0-9]+);([0-9]+)R" // x (line);y(column)
	setCursorPosition    string = "\x1b[%d;%dH"
	scrollUp             string = "\x1b[%dS"
	scrollDown           string = "\x1b[%dT"
	saveCurPos           string = "\x1b[s"
	restoreCurPos        string = "\x1b[u"
	hideCursor           string = "\x1b[?25l"
	showCursor           string = "\x1b[?25h"
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

type ClearMode interface {
	getModeInt() int
}

type clearMode int

func (c clearMode) getModeInt() int {
	return int(c)
}

const (
	clearScreen string = "\x1b[%dJ"

	ClearAfterCursor  clearMode = 0 // 0 - clears all from cursors position to the end of terminal
	ClearBeforeCursor clearMode = 1 // 1 - clears all from cursors position to the start of terminal
	ClearAll          clearMode = 2 // 2 - clears all
)

type ClearLineMode interface {
	getLineModeInt() int
}

type clearLineMode int

func (c clearLineMode) getLineModeInt() int {
	return int(c)
}

const (
	clearLine string = "\x1b[%dK"

	ClearLineAfterCursor  clearMode = 0 // 0 - clears all from cursors position to the end of Line
	ClearLineBeforeCursor clearMode = 1 // 1 - clears all from cursors position to the start of line
	ClearAllLine          clearMode = 2 // 2 - clears all line
)
