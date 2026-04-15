//go:build !windows

package platform

import (
	"os/exec"
)

func HideCmdWindow(cmd *exec.Cmd) {
	// Do nothing on non-Windows systems.
}
