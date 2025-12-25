package commands

import (
	"fmt"
	"os"

	"github.com/SteliosSpanos/mygit/pkg/objects"
	"github.com/SteliosSpanos/mygit/pkg/storage"
)

func CatFile(hash string) error {
	gitDir, err := FindGitDir()
	if err != nil {
		return err
	}

	obj, err := storage.LoadObject(gitDir, hash)
	if err != nil {
		return fmt.Errorf("failed to load object %s: %w", hash, err)
	}

	switch obj.Type() {
	case objects.BlobObject:
		blob := obj.(*objects.Blob)
		os.Stdout.Write(blob.Data)
	case objects.TreeObject:
		tree := obj.(*objects.Tree)
		PrintTree(tree)
	case objects.CommitObject:
		commit := obj.(*objects.Commit)
		PrintCommit(commit)
	default:
		return fmt.Errorf("unknown object type: %s", obj.Type())
	}

	return nil
}

func PrintTree(tree *objects.Tree) {
	for _, entry := range tree.Entries {
		var objType string

		switch entry.Mode {
		case "040000":
			objType = "tree"
		case "100644":
			objType = "blob"
		case "100755":
			objType = "blob"
		case "120000":
			objType = "blob"
		default:
			objType = "blob"
		}

		fmt.Printf("%s %s %s\t%s\n", entry.Mode, objType, entry.Hash, entry.Name)
	}
}

func PrintCommit(commit *objects.Commit) {
	fmt.Printf("tree %s\n", commit.Tree)

	for _, parent := range commit.Parents {
		fmt.Printf("parent %s\n", parent)
	}

	fmt.Printf("author %s\n", commit.Author)
	fmt.Printf("committer %s\n", commit.Committer)
	fmt.Printf("\n%s\n", commit.Message)
}
