package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	libDir := os.Getenv("LIB_DIR")
	extLibDir := os.Getenv("EXTERNAL_LIB_DIR")
	confDir := os.Getenv("CONF_DIR")
	extConfDir := os.Getenv("EXTERNAL_CONF_DIR")

	fmt.Printf("Lib directory is: '%s'\n", libDir)
	fmt.Printf("External lib directory is: '%s'\n", extLibDir)
	fmt.Printf("Config directory is: '%s'\n", confDir)
	fmt.Printf("External Config directory is: '%s'\n", extConfDir)

	_ = os.MkdirAll(extLibDir, os.ModePerm)
	_ = os.MkdirAll(extLibDir, os.ModePerm)

	_, err := os.Stat(libDir)
	if err == nil {
		err = copyDir(libDir, extLibDir)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("Lib directory '%s' does not exist - no files to copy\n", libDir)
	}

	_, err = os.Stat(confDir)
	if err == nil {
		err = copyDir(confDir, extConfDir)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("Config directory '%s' does not exist - no files to copy\n", confDir)
	}
}

func copyDir(src string, dst string) error {
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
			if err = copyDir(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcfp, dstfp); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// File copies a single file from src to dst
func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	fmt.Printf("Copied file %s to %s\n", src, dst)
	return os.Chmod(dst, srcinfo.Mode())
}
