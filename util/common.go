package util

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var CONFIG_PATH = ""
var DATA_DIR = ""

const (
	MASTER_URL = "https://raw.githubusercontent.com/runscripts/run/master/"
)

func SetConfigPath() {
	DATA_DIR = "/usr/local/run"
}
func SetDataDir() {
	DATA_DIR = "/usr/local/run"
}
func IsRunInstalled() bool {
	return FileExists(CONFIG_PATH) && FileExists(DATA_DIR)
}
func LogError(format string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, format, args...)
	} else {
		fmt.Fprintf(os.Stderr, format)
	}
}
func LogInfo(format string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Printf(format, args...)
	} else {
		fmt.Printf(format)
	}
}
func Errorf(format string, args ...interface{}) error {
	if len(args) > 0 {
		return fmt.Errorf(format, args...)
	} else {
		return fmt.Errorf(format)
	}
}
func ExitError(err error) {
	LogError("%v\n", err)
	os.Exit(1)
}

// Convert string into hash string.
func StrToSha1(str string) string {
	sum := [20]byte(sha1.Sum([]byte(str)))
	return hex.EncodeToString(sum[:])
}
func Exec(arg []string) {
	path, err := exec.LookPath(arg[0])
	if err != nil {
		ExitError(err)
	}
	err = syscall.Exec(path, arg, os.Environ())
	if err != nil {
		ExitError(err)
	}
}
func FileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
