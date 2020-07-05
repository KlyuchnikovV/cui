#include <stdlib.h>
#include <termios.h>
#include <unistd.h>

struct termios terminalState;

void disableRawMode() {
  tcsetattr(STDIN_FILENO, TCSAFLUSH, &terminalState);
}

void enableRawMode() {
  tcgetattr(STDIN_FILENO, &terminalState);
  atexit(disableRawMode);
  struct termios raw = terminalState;
  raw.c_lflag &= ~(ECHO | ICANON);
  tcsetattr(STDIN_FILENO, TCSAFLUSH, &raw);
}