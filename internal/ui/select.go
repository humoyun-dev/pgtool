package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SelectOne renders a numbered list and returns the chosen option.
func SelectOne(label string, options []string) (string, error) {
	fmt.Println(label)
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt)
	}
	fmt.Print("Select an option [number]: ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	idx, err := strconv.Atoi(line)
	if err != nil || idx < 1 || idx > len(options) {
		return "", fmt.Errorf("invalid selection")
	}
	return options[idx-1], nil
}

// SelectOneOrSkip renders a numbered list and returns the chosen option or the default when empty.
func SelectOneOrSkip(label string, options []string, def string) (string, error) {
	fmt.Println(label)
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt)
	}
	if def != "" {
		fmt.Printf("Select an option [number, empty = %s]: ", def)
	} else {
		fmt.Print("Select an option [number, empty = skip]: ")
	}

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return def, nil
	}
	idx, err := strconv.Atoi(line)
	if err != nil || idx < 1 || idx > len(options) {
		return "", fmt.Errorf("invalid selection")
	}
	return options[idx-1], nil
}
