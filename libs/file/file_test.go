package file_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mauricetjmurphy/ms-common/libs/file"
)

const (
	_testTxt = `
    	dummy testing data
	`
)

func TestRead(t *testing.T) {
	//Given
	var (
		path     = filepath.Join(t.TempDir(), "test_data")
		textFile = filepath.Join(path, "test.txt")
		data     = []byte(_testTxt)
	)
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0700); err != nil {
		t.Error(err)
	}
	if err := ioutil.WriteFile(textFile, data, 0666); err != nil {
		t.Error(err)
	}

	//When & Then
	testRead(t, textFile, data)
}

func testRead(t *testing.T, path string, data []byte) {
	t.Log(path)

	result, err := file.Read(path)
	if err != nil {
		t.Error(err)
	}
	if string(result) != string(data) {
		t.Errorf("no expected: %s, but got: %s", result, data)
	}
}
