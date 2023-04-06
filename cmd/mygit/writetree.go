package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteTreeCmd() error {
	dirname := os.Args[len(os.Args)-1]
	if len(os.Args) < 3 {
		dirname = "."
	}

	objectname, err := writeTree(dirname, true)
	if err != nil {
		return err
	}

	fmt.Println(objectname)

	return nil
}

func writeTree(dir string, write bool) (objectname string, err error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	tree := bytes.NewBuffer(nil)
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if file.Name()[0] == '.' {
				continue
			}

			objectname, err := writeTree(path, write)
			if err != nil {
				return "", nil
			}
			hexName, err := hex.DecodeString(objectname)
			if err != nil {
				return "", nil
			}
			fmt.Fprintf(tree, "40000 %s\x00%s", file.Name(), hexName)
			continue
		}

		objectname, err := hashObject(path, false)
		if err != nil {
			return "", err
		}
		hexName, err := hex.DecodeString(objectname)
		if err != nil {
			return "", nil
		}
		fmt.Fprintf(tree, "100644 %s\x00%s", file.Name(), hexName)
	}
	prefix := fmt.Sprintf("tree %d\x00", tree.Len())
	data := append([]byte(prefix), tree.Bytes()...)

	hasher := sha1.New()
	hasher.Write(data)
	objectname = hex.EncodeToString(hasher.Sum(nil))

	if write {
		err := WriteObject(objectname, data)
		if err != nil {
			return objectname, err
		}
	}

	return objectname, nil
}
