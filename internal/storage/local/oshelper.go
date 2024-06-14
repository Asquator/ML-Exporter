package local

import (
	"fmt"
	"os"
)

func openDir(path string) (*os.File, error) {
	stat, err := os.Stat(path)

	if err != nil {
		return nil, fmt.Errorf("error occurred: %w", err)
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("storage path is not a directory: %s", path)
	}

	return os.Open(path)
}
