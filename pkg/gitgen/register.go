package gitgen

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

var (
	ShellFiles = []string{".zshrc", ".zprofile", ".bashrc", ".bash_profile"}

	ErrShellFileNotFound = errors.New("shell file not found")
)

func CheckFileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func FindShellFile() (string, error) {
	userHome, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	for _, fileName := range ShellFiles {
		filePath := path.Join(userHome, fileName)

		if CheckFileExists(filePath) {
			return filePath, nil
		}
	}

	return "", ErrShellFileNotFound
}

func UpdatePathLine(existingLine string, newPath string) string {
	existingPaths := strings.Split(
		strings.Trim(existingLine[12:], `"`),
		":",
	)

	newPaths := []string{}
	hasDollarPath := false
	for _, path := range existingPaths {
		if path == "$PATH" {
			hasDollarPath = true
			continue
		}

		if path == newPath {
			continue
		}

		newPaths = append(newPaths, path)
	}

	if hasDollarPath {
		newPaths = append(newPaths, "$PATH")
	}

	newPaths = append(newPaths, newPath)

	return fmt.Sprintf(`export PATH="%s"`, strings.Join(newPaths, ":"))
}

func CreatePathLine(newPath string) string {
	return fmt.Sprintf(`export PATH="$PATH:%s"`, newPath)
}

func ModifyShellFile(filePath string, newPath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	var pathLineIndex int = -1

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "export PATH=") {
			pathLineIndex = len(lines)
			lines = append(lines, line)
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if pathLineIndex != -1 {
		lines[pathLineIndex] = UpdatePathLine(lines[pathLineIndex], newPath)
	} else {
		lines = append(lines, CreatePathLine(newPath))
	}

	output := strings.Join(lines, "\n")

	if err := os.WriteFile(filePath, []byte(output), 0644); err != nil {
		return err
	}

	return nil
}

func RegisterToPath() error {
	workingDirectory, err := os.Getwd()
	if err != nil {
		return err
	}

	shellFile, err := FindShellFile()
	if err != nil {
		return err
	}

	err = ModifyShellFile(shellFile, workingDirectory)
	if err != nil {
		return err
	}

	return nil
}
