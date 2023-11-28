package llog

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func OpenFile(fname string, flag int, perm fs.FileMode) *os.File {
	i := strings.LastIndexByte(fname, os.PathSeparator)
	// var pre string
	if i != -1 {
		pre := fname[:i]
		if err := os.MkdirAll(pre, perm); err != nil {
			fmt.Println("dir: " + pre + " created failed")
			return nil
		}
	}
	f, err := os.OpenFile(fname, flag, perm)
	if err != nil {
		fmt.Println("file: " + fname + " created failed")
		return nil
	}
	return f
}
