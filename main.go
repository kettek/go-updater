package updater

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// DebugOutput prints the execution steps to stdout.
var DebugOutput = false

func debugPrintln(a ...interface{}) {
	if !DebugOutput {
		return
	}
	debugPrintln(a...)
}

// DelayTime is the amount of time to delay before deleting the original file and moving the new file into the original file's location.
var DelayTime = uint(2)

func runAfter(target string, seconds uint) error {
	debugPrintln("Running", target, "in", seconds, "seconds...")
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/C", fmt.Sprintf("timeout /t %d && del %s.old && %s", seconds, target, target))
		return cmd.Start()
	} else {
		cmd := exec.Command("sh", "-c", fmt.Sprintf("sleep %d && rm %s.old && %s", seconds, target, target))
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		return cmd.Start()
	}
}

// Update replaces target with source and an optional PID to kill. If source begins with http, it is downloaded to a temp folder and the created file is used as the source.
func Update(source string, target string, pid string) error {
	isTemp := false

	if strings.HasPrefix(source, "http") {
		debugPrintln("Downloading", source, "...")
		resp, err := http.Get(source)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		filename := source[strings.LastIndex(source, "/")+1:]

		// Pipe body to temp file.
		temp, err := os.CreateTemp("", filename)
		if err != nil {
			return err
		}
		// Copy response body to temp.
		_, err = temp.ReadFrom(resp.Body)
		if err != nil {
			return err
		}
		temp.Close()
		source = temp.Name()
		isTemp = true
		debugPrintln("...ok")
	}

	// Rename ourself.
	debugPrintln("Renaming", target, "to", target+".old", "...")
	if err := os.Rename(target, target+".old"); err != nil {
		return err
	}
	debugPrintln("...ok")

	// Open target file for writing.
	debugPrintln("Creating", target, "...")
	targetFile, err := os.Create(target)
	if err != nil {
		panic(err)
	}
	debugPrintln("...ok")

	// Copy source to targetFile
	debugPrintln("Copying", source, "to", target, "...")
	sourceFile, err := os.Open(source)
	if err != nil {
		panic(err)
	}
	_, err = sourceFile.WriteTo(targetFile)
	if err != nil {
		panic(err)
	}
	sourceFile.Close()
	targetFile.Close()
	debugPrintln("...ok")

	// Mark targetFile as executable.
	debugPrintln("Marking", target, "as executable...")
	if err := os.Chmod(target, 0755); err != nil {
		return err
	}
	debugPrintln("...ok")

	if isTemp {
		debugPrintln("Removing", source, "...")
		if err := os.Remove(source); err != nil {
			debugPrintln(err)
		} else {
			debugPrintln("...ok")
		}
	}

	if err := runAfter(target, DelayTime); err != nil {
		debugPrintln(err)
	} else {
		debugPrintln("...ok")
	}

	// If PID is found, kill it.
	if pid != "" {
		debugPrintln("Killing PID", pid, "...")
		npid, err := strconv.ParseInt(pid, 10, 64)
		if err != nil {
			return err
		}
		p, err := os.FindProcess(int(npid))
		if err != nil {
			return err
		}
		err = p.Kill()
		if err != nil {
			return err
		}
		// Check if the process still exists.
		if err := p.Signal(syscall.Signal(0)); err != nil {
			return errors.New("process is still running")
		}
		debugPrintln("...ok")
	}

	return nil
}
