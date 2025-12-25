package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SteliosSpanos/mygit/pkg/objects"
	"github.com/SteliosSpanos/mygit/pkg/storage"
)

func HashObject(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filepath, err)
	}

	blob := objects.NewBlob(data)

	gitDir, err := FindGitDir()
	if err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}

	hash, err := storage.WriteObject(gitDir, blob)
	if err != nil {
		return fmt.Errorf("failed to write object: %w", err)
	}

	fmt.Println(hash)
	return nil
}

func FindGitDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get curret directory: %w", err)
	}

	for {
		gitDir := filepath.Join(dir, ".git")

		info, err := os.Stat(gitDir)
		if err == nil && info.IsDir() {
			return gitDir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not a git repository")
		}

		dir = parent
	}
}
