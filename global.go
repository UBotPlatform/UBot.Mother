package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/nicksnyder/basen"
)

var ExeSuffix = func() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}()

var PathExt = func() []string {
	if runtime.GOOS == "windows" {
		return []string{".bat", ".cmd", ".exe", ".com"}
	}
	return []string{".sh", ""}
}()

func LogFilePath(name string) string {
	return filepath.Join(LogFolder, name+".log")
}

func GetUBotAddr(scheme string, endpoint string) string {
	var addr string
	if Config.UBot.Address == "" {
		addr = ":5000"
	} else {
		addr = Config.UBot.Address
	}
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	return fmt.Sprintf("%s://%s%s", scheme, addr, endpoint)
}

func NewToken() string {
	var u [16]byte
	_, _ = io.ReadFull(rand.Reader, u[:])
	return basen.Base62Encoding.EncodeToString(u[:])
}

func ExitCmd(cmd *exec.Cmd) {
	if cmd == nil {
		return
	}
	process := cmd.Process
	if process == nil {
		return
	}
	if runtime.GOOS != "windows" {
		_ = process.Signal(syscall.SIGINT)
		if cmd.ProcessState == nil {
			go cmd.Wait() //nolint:errcheck
		}
		const timeout = 5 * time.Second
		const checkInterval = 100 * time.Millisecond
		const checkTimes = int(timeout / checkInterval)
		for i := 0; i < checkTimes; i++ {
			time.Sleep(checkInterval)
			if cmd.ProcessState.Exited() {
				return
			}
		}
	}
	_ = process.Kill()
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
