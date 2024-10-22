package updater

import (
	"fmt"
	"os/exec"
)

func runAfter(target string, seconds uint) error {
	fmt.Println("we in windows")
	cmd := exec.Command("cmd", "/C", fmt.Sprintf("start /b timeout /t %d && del %s.old && %s", seconds, target, target))
	return cmd.Start()
}
