package repository
import (
	"fmt"
	"os"
	"path/filepath"
)


const GitDir = ".git"


type Repository struct {
	WorkTree string
	GitDir string
}


func Init(path string) (*Repository, error) {
	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get current directory: %w", err)
		}
	}

	gitDir := filepath.Join(path, GitDir)

	if _, err := os.Stat(gitDir); err == nil {
		return nil, fmt.Errorf("repository already exists at %s", gitDir)
	}


	dirs := []string{
		gitDir,
		filepath.Join(gitDir, "objects"),
		filepath.Join(gitDir, "refs"),
		filepath.Join(gitDir, "refs", "heads"),
		filepath.Join(gitDir, "refs", "tags"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	headPath := filepath.Join(gitDir, "HEAD")
	headContent :=  "ref: refs/heads/main\n"

	if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to create HEAD: %w", err)
	}

	repo := &Repository{
		WorkTree: path,
		GitDir: gitDir,
	}

	return repo, nil
}
