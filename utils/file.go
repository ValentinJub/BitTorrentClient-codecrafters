package utils

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func ReadFile(file string) (*bytes.Buffer, error) {
	fileHandle, err := os.Open(file)
	if err != nil {
		return new(bytes.Buffer), fmt.Errorf("unable to open %s\nError: %s", file, err)
	}
	defer fileHandle.Close()
	// Put the file data in a buffer we can read from
	b := new(bytes.Buffer)
	_, err = io.Copy(b, fileHandle)
	if err != nil {
		return new(bytes.Buffer), fmt.Errorf("error while reading from the file: %s", err)
	}
	return b, nil
}

func SHA1Hash(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
