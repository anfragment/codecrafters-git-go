package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func CatFile(objectname string, writer io.Writer) error {
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
	if _, err := io.Copy(&contents, zlibReader); err != nil {
		return err
	}
	contents.ReadBytes('\x00')

	writer.Write(contents.Bytes())

	return nil
}
