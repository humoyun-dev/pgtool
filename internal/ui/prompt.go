package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// IsTerminal reports whether stdin is a terminal (interactive).
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// Prompt prints label and returns trimmed input from stdin.
func Prompt(label string) (string, error) {
	fmt.Fprint(os.Stdout, label)
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// PromptDefault asks with a default value shown in brackets; empty input returns def.
func PromptDefault(label, def string) (string, error) {
	label = strings.TrimSpace(label)
	var prompt string
	if def != "" {
		prompt = fmt.Sprintf("%s [%s]: ", label, def)
	} else {
		prompt = label
		if !strings.HasSuffix(prompt, ": ") {
			prompt += ": "
		}
	}
	answer, err := Prompt(prompt)
	if err != nil {
		return "", err
	}
	if answer == "" {
		return def, nil
	}
	return answer, nil
}

// PromptConfirm asks a yes/no style question with default.
// If defTrue is true, empty input counts as yes.
func PromptConfirm(label string, defTrue bool) (bool, error) {
	label = strings.TrimSpace(label)
	def := "N"
	if defTrue {
		def = "Y"
	}
	prompt := fmt.Sprintf("%s [%s]: ", label, def)
	answer, err := Prompt(prompt)
	if err != nil {
		return false, err
	}
	if answer == "" {
		return defTrue, nil
	}
	switch strings.ToLower(strings.TrimSpace(answer)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

// PromptPassword hides input while typing.
func PromptPassword(label string) (string, error) {
	fmt.Fprint(os.Stdout, label)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stdout)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
