package commands

import (
	"fmt"
	"os"
	"os/user"

	"github.com/SteliosSpanos/mygit/pkg/index"
	"github.com/SteliosSpanos/mygit/pkg/objects"
	"github.com/SteliosSpanos/mygit/pkg/refs"
	"github.com/SteliosSpanos/mygit/pkg/storage"
	"github.com/SteliosSpanos/mygit/pkg/tree"
)

func Commit(message string) error {
	gitDir, err := FindGitDir()
	if err != nil {
		return err
	}

	idx, err := index.ReadIndex(gitDir)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	if len(idx.Entries) == 0 {
		return fmt.Errorf("nothing to commit (index is empty)")
	}

	treeHash, err := tree.BuildTreeFromIndex(gitDir, idx)
	if err != nil {
		return fmt.Errorf("failed to build tree: %w", err)
	}

	author := getAuthor()

	commit := objects.NewCommit(treeHash, author, message)

	currentBranch, err := refs.GetCurrentBranch(gitDir)
	if err != nil {
		return fmt.Errorf("failed to get current branch: %w", err)
	}

	parentHash, err := refs.ReadRef(gitDir, currentBranch)
	if err != nil {
		return fmt.Errorf("failed to read branch: %w", err)
	}

	if parentHash != "" {
		commit.AddParent(parentHash)
	}

	commitHash, err := storage.WriteObject(gitDir, commit)
	if err != nil {
		return fmt.Errorf("failed to write commit: %w", err)
	}

	if err := refs.WriteRef(gitDir, currentBranch, commitHash); err != nil {
		return fmt.Errorf("failed to update branch: %w", err)
	}

	branchName := currentBranch[len("refs/heads/"):]
	if parentHash == "" {
		fmt.Printf("[%s (root-commit) %s] %s\n", branchName, commitHash[:7], message)
	} else {
		fmt.Printf("[%s %s] %s\n", branchName, commitHash[:7], message)
	}

	return nil
}

func getAuthor() string {
	name := os.Getenv("GIT_AUTHOR_NAME")
	email := os.Getenv("GIT_AUTHOR_EMAIL")

	if name == "" || email == "" {
		currentUser, err := user.Current()
		if err == nil {
			name = currentUser.Username
		} else {
			name = "Unknown"
		}

		email = name + "@localhost"
	}

	return fmt.Sprintf("%s <%s>", name, email)
}
