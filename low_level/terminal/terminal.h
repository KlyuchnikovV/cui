#include <sys/ioctl.h>
#include <unistd.h>

struct Size
{
    int width;
    int height;
};

int getTerminalWidth();
int getTerminalHeight();
struct Size getTerminalSize();
