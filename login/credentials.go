package login

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
)

func GetSavedLogin() (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	creds := filepath.Join(home, ".cap-credentials")
	file, err := os.Open(creds)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}
	if len(lines) < 2 {
		return "", "", errors.New("Could not read .cap-credentials")
	}
	return lines[0], lines[1], nil
}
