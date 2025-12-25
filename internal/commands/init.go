package commands
import (
	"fmt"
	"github.com/SteliosSpanos/mygit/pkg/repository"
)


func Init() error {
	repo, err := repository.Init("")
	if err != nil {
		return err
	}

	fmt.Printf("Initialized empty repository in %s\n", repo.GitDir)
	return nil
}
