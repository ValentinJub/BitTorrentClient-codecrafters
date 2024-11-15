package utils

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
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

func WriteFile(file string, data []byte) error {
	fileHandle, err := os.Create(file)
	if err != nil {
		return fmt.Errorf("unable to create %s\nError: %s", file, err)
	}
	defer fileHandle.Close()
	_, err = fileHandle.Write(data)
	if err != nil {
		return fmt.Errorf("error while writing to the file: %s", err)
	}
	return nil
}

func SHA1Hash(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return string(hash.Sum(nil))
}

func GeneratePeerID() string {
	return "-PC0001-" + RandStringBytes(12)
}

func RandStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
