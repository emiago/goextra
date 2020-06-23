package osx

import (
	"io/ioutil"
	"os"
	"syscall"
)

//Check does file exists, it also check for syscall errors
func CheckFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == os.ErrNotExist {
		return false
	}

	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
		return false
	}
	return true
}

//This compares two files are same by size and content
func CompareFilesAreSame(afile, bfile string) bool {
	res := CheckFileExists(afile) && CheckFileExists(bfile)
	if !res {
		return false
	}

	afi, err := os.Stat(afile)
	if err != nil {
		return false
	}

	bfi, err := os.Stat(bfile)
	if err != nil {
		return false
	}

	if afi.IsDir() {
		return bfi.IsDir()
	}

	if bfi.IsDir() {
		return afi.IsDir()
	}

	res = afi.Size() == bfi.Size()
	if !res {
		return false
	}

	adata, _ := ioutil.ReadFile(afile)
	bdata, _ := ioutil.ReadFile(bfile)
	return string(adata) == string(bdata)
}
