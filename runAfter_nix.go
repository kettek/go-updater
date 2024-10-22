//go:build !windows

package updater

import (
	"fmt"
	"os/exec"
	"syscall"
)

func runAfter(target string, seconds uint) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("sleep %d && rm %s.old && %s", seconds, target, target))
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd.Start()
}
