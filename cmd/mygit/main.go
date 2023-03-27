package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	p = flag.String("p", "", "")
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	command := os.Args[1]
	os.Args = append(os.Args[:1], os.Args[2:]...)
	flag.Parse()

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
		if p == nil {
			log.Fatal("cat-file: -p parameter required")
		}
		if len(*p) != 40 {
			log.Fatal("cat-file: -p has to be 40 characters long")
		}

		dir, fileName := (*p)[:2], (*p)[2:]
		file, err := os.Open(fmt.Sprintf(".git/objects/%s/%s", dir, fileName))
		if err != nil {
			log.Fatal("cat-file: unable to open file")
		}
		defer file.Close()

		zlibReader, err := zlib.NewReader(file)
		if err != nil {
			log.Fatal("cat-file: unable to initialize zlib reader")
		}
		defer zlibReader.Close()

		var contents bytes.Buffer
		if _, err := io.Copy(&contents, zlibReader); err != nil {
			log.Fatal("cat-file: error")
		}

		contents.ReadBytes('\x00')

		fmt.Print(contents.String())

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
