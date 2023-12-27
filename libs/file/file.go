package file

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Read utility function open and read the named file.
func Read(path string) ([]byte, error) {
	in, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := in.Close(); err != nil {
			log.Printf("error on closing file: %s\n", err)
		}
	}()
	return ioutil.ReadAll(in)
}
