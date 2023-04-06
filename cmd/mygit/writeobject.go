package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
)

// WriteObject compressed the data and writes it into the .git/objects directory
func WriteObject(objectname string, data []byte) error {
	compressed := bytes.NewBuffer(nil)
	zlibWriter := zlib.NewWriter(compressed)
	if _, err := zlibWriter.Write(data); err != nil {
		return err
	}
	if err := zlibWriter.Close(); err != nil {
		return err
	}

	dirname, filename := objectname[:2], objectname[2:]
	err := os.MkdirAll(fmt.Sprintf(".git/objects/%s", dirname), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf(".git/objects/%s/%s", dirname, filename), compressed.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
