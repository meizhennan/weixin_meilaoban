package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteFile(buf []byte, fileName string) (string, error) {
	tmpPath := "./tmp"
	err := os.MkdirAll(tmpPath, 0700)
	if err != nil {
		return "", err
	}

	file := filepath.Join(tmpPath, fileName)
	err = ioutil.WriteFile(file, buf, 0600)
	if err != nil {
		return "", err
	}
	return file, err
}
