package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func CatFile(objectname string) error {
	dir, filename := objectname[:2], objectname[2:]
	file, err := os.Open(fmt.Sprintf(".git/objects/%s/%s", dir, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return err
	}
	defer zlibReader.Close()

	var contents bytes.Buffer
	_, err = io.Copy(&contents, zlibReader)
	if err != nil {
		return err
	}
	contents.ReadBytes('\x00')

	fmt.Print(string(contents.Bytes()))

	return nil
}
