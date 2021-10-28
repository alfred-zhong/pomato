package pomato

/*
#include <termios.h>
#include <unistd.h>

void flushtty() {
tcflush(0, TCIFLUSH);
}
*/
import "C"

func flushStdin() {
	C.flushtty()
}
