package main

import (
	"fmt"
	"os"

	"github.com/SteliosSpanos/mygit/internal/commands"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: mygit <command> [<args>]")
		fmt.Println("Commands:")
		fmt.Println("   init   Initialize a new repository")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		err := commands.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "hash-object":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: mygit hash-object <file>\n")
			os.Exit(1)
		}

		filePath := os.Args[2]
		if err := commands.HashObject(filePath); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
