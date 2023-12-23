package flock

import (
	"os"
	"syscall"
)

var lockFile *os.File

func Flock(path string) error {
	return fcntlFlock(syscall.F_WRLCK, path)
}
func Funlock(path string) error {
	err := fcntlFlock(syscall.F_WRLCK, path)
	if err != nil {
		return err
	} else {
		return lockFile.Close()
	}
}
func fcntlFlock(lockType int16, path ...string) error {
	var err error
	if lockType != syscall.F_UNLCK {
		mode := syscall.O_CREAT | syscall.O_WRONLY
		mask := syscall.Umask(0)
		lockFile, err = os.OpenFile(path[0], mode, 0666)
		syscall.Umask(mask)
		if err != nil {
			return err
		}
	}
	lock := syscall.Flock_t{
		Start:  0,
		Len:    1,
		Type:   lockType,
		Whence: int16(os.SEEK_SET),
	}
	return syscall.FcntlFlock(lockFile.Fd(), syscall.F_SETLK, &lock)
}
