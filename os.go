package shell

import (
	"os"
	"runtime"
)

// IsDebian 判断是否是 Debian 系统
func IsDebian() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	_, err := os.Stat("/etc/debian_version")
	return err == nil
}

// IsRHEL 判断是否是 RHEL 系统
func IsRHEL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	_, err := os.Stat("/etc/redhat-release")
	return err == nil
}
