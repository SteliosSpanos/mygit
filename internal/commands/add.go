package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SteliosSpanos/mygit/pkg/index"
	"github.com/SteliosSpanos/mygit/pkg/objects"
	"github.com/SteliosSpanos/mygit/pkg/storage"
)

func Add(filePath string) error {
	gitDir, err := FindGitDir()
	if err != nil {
		return err
	}

	repoRoot := filepath.Dir(gitDir)

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	relPath, err := filepath.Rel(repoRoot, absPath)
	if err != nil {
		return fmt.Errorf("file is outside repository: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	blob := objects.NewBlob(data)
	hash, err := storage.WriteObject(gitDir, blob)
	if err != nil {
		return fmt.Errorf("failed to write blob: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	mode := "100644"
	if info.Mode()&0111 != 0 {
		mode = "100755"
	}

	idx, err := index.ReadIndex(gitDir)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	idx.Add(mode, hash, relPath)

	if err := index.WriteIndex(gitDir, idx); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	fmt.Printf("Added '%s' to staging area\n", relPath)
	return nil
}
