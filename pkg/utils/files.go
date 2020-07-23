/*
 * Copyright (c) 2019, Oracle and/or its affiliates.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type FileFilter func(string) bool

var alwaysFilter FileFilter = func(string) bool { return true }

func AlwaysFilter() FileFilter {
	return alwaysFilter
}

func CopyDir(src string, dst string, filter FileFilter) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	fmt.Printf("Copying directory %s to %s\n", src, dst)

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp, filter); err != nil {
				fmt.Println(err)
			}
		} else if filter(srcfp) {
			if err = CopyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// File copies a single file from src to dst
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer closeFile(srcfd)

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer closeFile(dstfd)

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	fmt.Printf("Copied file %s to %s\n", src, dst)
	return os.Chmod(dst, srcinfo.Mode())
}

func closeFile(f *os.File) {
	if f != nil {
		_ = f.Close()
	}
}
