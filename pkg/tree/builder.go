package tree

import (
	"fmt"
	"strings"

	"github.com/SteliosSpanos/mygit/pkg/index"
	"github.com/SteliosSpanos/mygit/pkg/objects"
	"github.com/SteliosSpanos/mygit/pkg/storage"
)

func BuildTreeFromIndex(gitDir string, idx *index.Index) (string, error) {
	if len(idx.Entries) == 0 {
		return "", fmt.Errorf("nothing to commit (index is empty)")
	}

	rootTree := buildTree(idx.Entries, "")

	hash, err := storeTree(gitDir, rootTree)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func buildTree(entries []index.Entry, prefix string) *objects.Tree {
	tree := objects.NewTree()
	subDirs := make(map[string][]index.Entry)

	for _, entry := range entries {
		relPath := entry.Path
		if prefix != "" {
			relPath = strings.TrimPrefix(relPath, prefix+"/")
		}

		parts := strings.Split(relPath, "/")
		if len(parts) == 1 {
			tree.AddEntry(entry.Mode, parts[0], entry.Hash)
		} else {
			subDir := parts[0]
			if subDirs[subDir] == nil {
				subDirs[subDir] = make([]index.Entry, 0)
			}
			subDirs[subDir] = append(subDirs[subDir], entry)
		}
	}

	//For now we only support all files in root
	return tree
}

func storeTree(gitDir string, tree *objects.Tree) (string, error) {
	hash, err := storage.WriteObject(gitDir, tree)
	if err != nil {
		return "", fmt.Errorf("failed to store tree: %w", err)
	}

	return hash, nil
}
