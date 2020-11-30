package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func NewTempFile(pattern string) (file *os.File, err error) {
	filenamePattern := filepath.Clean(fmt.Sprint(pattern))
	file, err = ioutil.TempFile(os.TempDir(), filenamePattern)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		err = nil
	}
	return
}
