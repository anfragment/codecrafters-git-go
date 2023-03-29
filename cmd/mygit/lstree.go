package main

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func LsTree(objectname string, nameonly bool) (err error) {
	dir, filename := objectname[:2], objectname[2:]
	file, err := os.Open(fmt.Sprintf(".git/objects/%s/%s", dir, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	zLibReader, err := zlib.NewReader(file)
	if err != nil {
		return err
	}
	defer zLibReader.Close()

	var contents bytes.Buffer
	_, err = io.Copy(&contents, zLibReader)
	if err != nil {
		return err
	}
	contents.ReadBytes('\x00')

	for {
		mode, err := contents.ReadBytes(' ')
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if !nameonly {
			fmt.Printf("%s\t", string(mode))
		}

		objectname, err := contents.ReadBytes('\x00')
		if err != nil {
			return err
		}
		fmt.Printf("%s", objectname[:len(objectname)-1])
		if !nameonly {
			fmt.Print(" ")
		}

		sha := make([]byte, 20)
		_, err = contents.Read(sha)
		if err != nil {
			return err
		}
		if !nameonly {
			fmt.Print(hex.EncodeToString(sha))
		}
		fmt.Print("\n")
	}

	return nil
}
