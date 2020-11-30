package utils

import (
	"path/filepath"
)

// ExpandDir will expand input to absolute path
func ExpandDir(dir string) (wd string, err error) {
	if !filepath.IsAbs(wd) {
		wd, err = filepath.Abs(dir)
	} else {
		wd = dir
	}
	return
}
