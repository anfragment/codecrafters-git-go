package main

import (
	"fmt"
	"os"
)

var ()

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		err := InitCmd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "init: %v\n", err)
			os.Exit(1)
		}

	case "cat-file":
		err := CatFileCmd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat-file: %v\n", err)
			os.Exit(1)
		}

	case "hash-object":
		err := HashObjectCmd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "hash-object: %v\n", err)
			os.Exit(1)
		}

	case "ls-tree":
		err := LsTreeCmd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls-tree: %v\n", err)
			os.Exit(1)
		}

	case "write-tree":
		err := WriteTreeCmd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "write-tree: %v\n", err)
			os.Exit(1)
		}

	case "commit-tree":
		err := CommitTreeCmd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "commit-tree: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
