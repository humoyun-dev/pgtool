package sys

import (
	"os"
	"os/exec"
)

func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func RunCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
