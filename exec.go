package shell

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Exec(shell string, args ...any) (string, error) {
	if err := checkArgs(args...); err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	_ = os.Setenv("LC_ALL", "C")
	cmd = exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", errors.New(strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), err
}

// ExecAsync 异步执行 shell 命令
func ExecAsync(shell string, args ...any) error {
	if err := checkArgs(args...); err != nil {
		return err
	}

	var cmd *exec.Cmd
	_ = os.Setenv("LC_ALL", "C")
	cmd = exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

// ExecWithTimeout 执行 shell 命令并设置超时时间
func ExecWithTimeout(timeout time.Duration, shell string, args ...any) (string, error) {
	if err := checkArgs(args...); err != nil {
		return "", err
	}

	var cmd *exec.Cmd
	_ = os.Setenv("LC_ALL", "C")
	cmd = exec.Command("bash", "-c", fmt.Sprintf(shell, args...))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return "", err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		_ = cmd.Process.Kill()
		return "", errors.New("timeout")
	case err = <-done:
		if err != nil {
			return "", errors.New(strings.TrimSpace(stderr.String()))
		}
	}

	return strings.TrimSpace(stdout.String()), err
}

// checkArgs 检查危险的参数
func checkArgs(args ...any) error {
	if len(args) == 0 {
		return nil
	}

	dangerous := []string{"&", "|", ";", "$", "'", `"`, "(", ")", "`", "\n", "\r", ">", "<", "{", "}", "[", "]", "\\"}
	for _, arg := range args {
		for _, char := range dangerous {
			if strings.Contains(arg.(string), char) {
				return errors.New("error args")
			}
		}
	}

	return errors.New("error args")
}
