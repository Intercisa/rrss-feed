package util

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/term"
)

// OpenFileUtil defines the system utility to use to open files
var (
	OpenFileUtil = "open"
	OpenUrlUtil  = []string{}
)

func OpenFile(path string) {
	if (strings.HasPrefix(path, "http://")) || (strings.HasPrefix(path, "https://")) {
		if len(OpenUrlUtil) > 0 {
			commands := append(OpenUrlUtil, path)
			cmd := exec.Command(commands[0], commands[1:]...)
			err := cmd.Start()
			if err != nil {
				return
			}
			return
		}

		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "linux":
			cmd = exec.Command("xdg-open", path)
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
		case "darwin":
			cmd = exec.Command("open", path)
		default:
			// for the BSDs
			cmd = exec.Command("xdg-open", path)
		}

		err := cmd.Start()
		if err != nil {
			return
		}
		return
	}

	filePath, _ := ExpandHomeDir(path)
	cmd := exec.Command(OpenFileUtil, filePath)
	ExecuteCommand(cmd)
}

func ExpandHomeDir(path string) (string, error) {
	if path == "" {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) > 1 && path[1] != '/' && path[1] != '\\' {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[1:]), nil
}

func ExecuteCommand(cmd *exec.Cmd) string {
	if cmd == nil {
		return ""
	}

	buf := &bytes.Buffer{}
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return err.Error()
	}

	return buf.String()
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "linux", "darwin":
		fmt.Print("\033[H\033[2J")
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Unsupported platform!")
	}
}

func GetTermWidth() int {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting terminal size: %v\n", err)
		return 150
	}
	return width
}

func Turnicate(input string, limit int) string {
	if len(input) < 1 {
		return input
	}

	turnicated := ""

	for i, c := range input {
		turnicated += string(c)

		if (i + 1) == limit {
			break
		}
	}
	return turnicated
}
