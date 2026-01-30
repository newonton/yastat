package lock

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	dirName  = ".config/yastat"
	lockName = "yastat.lock"
)

func Acquire() (*os.File, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, lockName)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("yastat already running")
		}
		return nil, err
	}

	fmt.Fprintf(f, "%d\n", os.Getpid())

	return f, nil
}

func Release(f *os.File) {
	if f == nil {
		return
	}
	path := f.Name()
	f.Close()
	os.Remove(path)
}
