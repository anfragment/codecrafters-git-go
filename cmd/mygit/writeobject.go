package main

import (
	"fmt"
	"os"
)

func WriteObject(objectname string, contents []byte) error {
	dirname, filename := objectname[:2], objectname[2:]
	err := os.MkdirAll(fmt.Sprintf(".git/objects/%s", dirname), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf(".git/objects/%s/%s", dirname, filename), contents, 0644)
	if err != nil {
		return err
	}

	return nil
}
