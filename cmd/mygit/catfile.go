package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"os"
)

func CatFileCmd() error {
	p := flag.Bool("p", false, "Pretty-print object's content")

	objectname := os.Args[len(os.Args)-1]

	os.Args = os.Args[1 : len(os.Args)-1]
	flag.Parse()
	if !*p {
		return fmt.Errorf("-p required")
	}

	err := catFile(objectname)
	if err != nil {
		return err
	}
	return nil
}

func catFile(objectname string) error {
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

	fmt.Print(contents.String())

	return nil
}
