package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	p = flag.String("p", "", "object which contents to print")
	w = flag.String("w", "", "file which contents to write to the objects database")
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

	case "hash-object":
		if w == nil {
			log.Fatal("hash-object: -w parameter required")
		}
		// read file once to both compute sha-1 sum and compress using zlib
		// with larger files, using os.Open might be better
		contents, err := ioutil.ReadFile(*w)
		if err != nil {
			log.Fatal("hash-object: unable to open file")
		}
		// https://alblue.bandlem.com/2011/08/git-tip-of-week-objects.html
		prefix := fmt.Sprintf("blob %d\x00", len(contents))
		data := append([]byte(prefix), contents...)

		hasher := sha1.New()
		hasher.Write(data)
		objectName := hex.EncodeToString(hasher.Sum(nil))

		dirName, fileName := objectName[:2], objectName[2:]
		if err := os.MkdirAll(fmt.Sprintf(".git/objects/%s", dirName), 0755); err != nil {
			log.Fatal("hash-object: unable to create parent object directory")
		}

		compressed := bytes.NewBuffer(nil)
		zlibWriter := zlib.NewWriter(compressed)
		if _, err := zlibWriter.Write(data); err != nil {
			log.Fatal("hash-object: failed to compress the file")
		}
		if err := zlibWriter.Close(); err != nil {
			log.Fatal("hash-object: failed to flush the compressor")
		}
		if err := os.WriteFile(fmt.Sprintf(".git/objects/%s/%s", dirName, fileName), compressed.Bytes(), 0644); err != nil {
			log.Fatal("hash-object: failed to write to object file")
		}
		fmt.Println(objectName)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
