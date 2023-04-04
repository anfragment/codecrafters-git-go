package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func WriteTree(dir string, write bool) (objectname string, err error) {
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

			objectname, err := WriteTree(path, write)
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

		objectname, err := HashObject(path, false)
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

	compressed := bytes.NewBuffer(make([]byte, 0))
	zlibWriter := zlib.NewWriter(compressed)
	if _, err := zlibWriter.Write(data); err != nil {
		return "", err
	}
	if err := zlibWriter.Close(); err != nil {
		return "", err
	}
	if write {
		err := WriteObject(objectname, compressed.Bytes())
		if err != nil {
			return objectname, err
		}
	}

	return objectname, nil
}
