// +build windows

package main

import (
	"syscall"
)

func setDeathsig(_ *syscall.SysProcAttr) {
}
