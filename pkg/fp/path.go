package fp

import (
	"os"
	"path/filepath"
)

func CreateNonExistingFolder(path string) error {
	path, _ = filepath.Abs(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0o700)
	} else if err != nil {
		return err
	}
	return nil
}
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
