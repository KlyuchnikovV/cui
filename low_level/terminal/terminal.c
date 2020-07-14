#include "terminal.h"

struct winsize loadTerminalSize() {
    struct winsize size;
    ioctl(STDOUT_FILENO, TIOCGWINSZ, &size);
    return size;
}

int getTerminalWidth() {
    return loadTerminalSize().ws_col;
}

int getTerminalHeight() {
    return loadTerminalSize().ws_row;
}

struct Size getTerminalSize() {
    struct winsize size = loadTerminalSize();
    struct Size result = {
       width: size.ws_col,
       height: size.ws_row,
    };
    return result;   
}