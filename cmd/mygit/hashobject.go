package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func HashObjectCmd() error {
	w := flag.Bool("w", false, "Actually write the object into the object database")

	filepath := os.Args[len(os.Args)-1]

	os.Args = os.Args[1 : len(os.Args)-1]
	flag.Parse()

	objectname, err := hashObject(filepath, *w)
	if err != nil {
		return err
	}

	fmt.Println(objectname)

	return nil
}

func hashObject(filepath string, write bool) (objectname string, err error) {
	// read file once to both compute sha-1 sum and compress using zlib
	// with larger files, using os.Open might be better
	contents, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	// https://alblue.bandlem.com/2011/08/git-tip-of-week-objects.html
	prefix := fmt.Sprintf("blob %d\x00", len(contents))
	data := append([]byte(prefix), contents...)

	hasher := sha1.New()
	hasher.Write(data)
	objectname = hex.EncodeToString(hasher.Sum(nil))

	if write {
		err := WriteObject(objectname, data)
		if err != nil {
			return objectname, nil
		}
	}

	return objectname, nil
}
