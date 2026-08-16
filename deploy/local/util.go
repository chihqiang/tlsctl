package local

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func CopyFile(path string, content []byte) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("fail to write to file: %w", err)
	}
	return nil
}
func ExecCommand(command, shellEnv string) (string, string, error) {
	var cmd *exec.Cmd
	switch shellEnv {
	case "sh":
		cmd = exec.Command("sh", "-c", command)
	case "cmd":
		cmd = exec.Command("cmd", "/C", command)
	case "powershell":
		cmd = exec.Command("powershell", "-Command", command)
	default:
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", command)
		} else {
			cmd = exec.Command("sh", "-c", command)
		}
	}
	stdoutBuf := bytes.NewBuffer(nil)
	cmd.Stdout = stdoutBuf
	stderrBuf := bytes.NewBuffer(nil)
	cmd.Stderr = stderrBuf
	err := cmd.Run()
	if err != nil {
		return stdoutBuf.String(), stderrBuf.String(), errors.Wrap(err, "failed to execute command")
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}
