package refs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadRef(gitDir, refName string) (string, error) {
	if refName == "HEAD" {
		headPath := filepath.Join(gitDir, "HEAD")
		data, err := os.ReadFile(headPath)
		if err != nil {
			return "", fmt.Errorf("failed to read HEAD: %w", err)
		}

		content := strings.TrimSpace(string(data))

		// HEAD can be "ref: refs/heads/main" or direct hash
		if strings.HasPrefix(content, "ref: ") {
			targetRef := strings.TrimPrefix(content, "ref: ")
			return ReadRef(gitDir, targetRef)
		}

		return content, nil
	}

	refPath := filepath.Join(gitDir, refName)
	data, err := os.ReadFile(refPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // Branch doesnt exist yet
		}

		return "", fmt.Errorf("failed to read ref %s: %w", refName, err)
	}

	return strings.TrimSpace(string(data)), nil
}

func WriteRef(gitDir, refName, commitHash string) error {
	refPath := filepath.Join(gitDir, refName)

	refDir := filepath.Dir(refPath)
	if err := os.MkdirAll(refDir, 0755); err != nil {
		return fmt.Errorf("failed to create ref directory: %w", err)
	}

	if err := os.WriteFile(refPath, []byte(commitHash+"\n"), 0644); err != nil {
		return fmt.Errorf("failed to write ref: %w", err)
	}

	return nil
}

func GetCurrentBranch(gitDir string) (string, error) {
	headPath := filepath.Join(gitDir, "HEAD")

	data, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %w", err)
	}

	content := strings.TrimSpace(string(data))

	if strings.HasPrefix(content, "ref: ") {
		return strings.TrimPrefix(content, "ref: "), nil
	}

	return "", fmt.Errorf("HEAD is detached (not pointing to a branch)")
}
