package pomato

import (
	"golang.org/x/sys/unix"
)

func flushStdin() {
	unix.Syscall(unix.SYS_IOCTL, 0, unix.TCFLSH, 0)
}
