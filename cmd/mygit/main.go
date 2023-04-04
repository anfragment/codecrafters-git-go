package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	p        = flag.Bool("p", false, "Pretty-print object's content")
	w        = flag.Bool("w", false, "Actually write the object into the object database")
	nameonly = flag.Bool("name-only", false, "List only filenames (instead of the \"long\" output), one per line.")
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/master\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")
	case "cat-file":
		objectname := os.Args[len(os.Args)-1]

		os.Args = os.Args[1 : len(os.Args)-1]
		flag.Parse()
		if !*p {
			fmt.Fprintf(os.Stderr, "cat-file: -p required\n")
			os.Exit(1)
		}

		err := CatFile(objectname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cat-file: %v\n", err)
		}

	case "hash-object":
		filepath := os.Args[len(os.Args)-1]

		os.Args = os.Args[1 : len(os.Args)-1]
		flag.Parse()

		objectname, err := HashObject(filepath, *w)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hash-object: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(objectname)

	case "ls-tree":
		objectname := os.Args[len(os.Args)-1]

		os.Args = os.Args[1 : len(os.Args)-1]
		flag.Parse()

		err := LsTree(objectname, *nameonly)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls-tree: %v\n", err)
			os.Exit(1)
		}
	case "write-tree":
		dirname := os.Args[len(os.Args)-1]
		if len(os.Args) < 3 {
			dirname = "."
		}

		objectname, err := WriteTree(dirname, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write-tree: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(objectname)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
