package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"time"
)

const author = "grumpydoggo <me@grumpydog.space>"

func CommitTreeCmd() error {
	p := flag.String("p", "", "Indicates the id of a parent commit object")
	m := flag.String("m", "", "A paragraph in the commit log message")

	treeSha := os.Args[2]

	os.Args = os.Args[2:len(os.Args)]
	flag.Parse()
	if *p == "" {
		return fmt.Errorf("-p is required")
	}
	if *m == "" {
		return fmt.Errorf("-m is required")
	}

	objectname, err := commitTree(treeSha, *p, *m)
	if err != nil {
		return err
	}
	fmt.Println(objectname)

	return nil
}

func commitTree(treeSha, commitSha, message string) (objectname string, err error) {
	commit := bytes.NewBuffer(nil)

	fmt.Fprintf(commit, "tree %s\n", treeSha)
	fmt.Fprintf(commit, "parent %s\n", commitSha)

	t := timestamp()
	fmt.Fprintf(commit, "author %s %s\n", author, t)
	fmt.Fprintf(commit, "committer %s %s\n\n", author, t)

	commit.WriteString(message)
	commit.WriteRune('\n')

	prefix := fmt.Sprintf("commit %d\x00", commit.Len())
	data := append([]byte(prefix), commit.Bytes()...)

	hasher := sha1.New()
	hasher.Write(data)
	objectname = hex.EncodeToString(hasher.Sum(nil))

	err = WriteObject(objectname, data)
	if err != nil {
		return objectname, err
	}

	return objectname, nil
}

func timestamp() string {
	now := time.Now()

	unix := now.Unix()

	_, offset := now.Zone()
	offsetH := abs(offset / 60)
	tz := fmt.Sprintf("%02d%02d", offsetH/60, offsetH%60)
	if offset < 0 {
		tz = "-" + tz
	}

	return fmt.Sprintf("%d %s", unix, tz)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
