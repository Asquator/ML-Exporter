package local

import (
	"fmt"
	"os"
)

type Storage struct {
	dirpath string
}

// New creates a new local storage instance.
// Returns an error if the storage cannot be created or opened.
func New(storagePath string) (*Storage, error) {
	// Check if the directory exists.
	if stat, err := os.Stat(storagePath); os.IsExist(err) && stat.IsDir() {
		// If the directory exists, return a new storage instance.
		return &Storage{dirpath: storagePath}, nil
	}

	// If the directory doesn't exist, create it.
	err := os.MkdirAll(storagePath, 0755)
	if err != nil {
		// Return an error if the directory cannot be created.
		return nil, fmt.Errorf("cannot create storage: %w", err)
	}

	// Return a new storage instance.
	return &Storage{dirpath: storagePath}, nil
}
