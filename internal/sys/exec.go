package sys

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CommandExists reports whether the binary is present in PATH.
func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// RunCommand logs and executes a command, wiring stdout/stderr through.
func RunCommand(label string, name string, args ...string) error {
	fmt.Printf("- %s: %s %s\n", label, name, strings.Join(args, " "))
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RemovePath removes a filesystem path if it exists, logging the action.
func RemovePath(label, path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("- %s: %s (not found, skipping)\n", label, path)
			return nil
		}
		return err
	}
	fmt.Printf("- %s: removing %s\n", label, path)
	return os.RemoveAll(path)
}

// HomeDir returns the current user's home directory or "." on failure.
func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return home
}
