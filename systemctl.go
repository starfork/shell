package shell

import (
	"errors"
	"os/exec"
	"strings"
)

// Status
func Status(name string) (bool, error) {
	output, err := Exec("systemctl status %s | grep Active | grep -v grep | awk '{print $2}'", name)
	return output == "active", err
}

// IsEnabled
func IsEnabled(name string) (bool, error) {
	cmd := exec.Command("systemctl", "is-enabled", name)
	output, _ := cmd.CombinedOutput()
	status := strings.TrimSpace(string(output))

	switch status {
	case "enabled":
		return true, nil
	case "disabled":
		return false, nil
	case "masked":
		return false, errors.New("service is masked")
	case "static":
		return false, errors.New("service is enabled static")
	case "indirect":
		return false, errors.New("service is enabled indirect")
	default:
		return false, errors.New("unknow service status")
	}
}

// Start
func Start(name string) error {
	_, err := Exec("systemctl start %s", name)
	return err
}

// Stop
func Stop(name string) error {
	_, err := Exec("systemctl stop %s", name)
	return err
}

// Restart
func Restart(name string) error {
	_, err := Exec("systemctl restart %s", name)
	return err
}

// Reload
func Reload(name string) error {
	_, err := Exec("systemctl reload %s", name)
	return err
}

// Enable
func Enable(name string) error {
	_, err := Exec("systemctl enable %s", name)
	return err
}

// Disable
func Disable(name string) error {
	_, err := Exec("systemctl disable %s", name)
	return err
}
