// +build darwin freebsd openbsd netbsd dragonfly
// +build !appengine

package logrus1

import "syscall"

const ioctlReadTermios = syscall.TIOCGETA

type Termios syscall.Termios
