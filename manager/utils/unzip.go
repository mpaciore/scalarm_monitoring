package utils

import (
	"archive/zip"
	"os"
	"io"
	"path/filepath"
)

func cloneZipItem(f *zip.File, dest string){
	// Create full directory path
	path := filepath.Join(dest, f.Name)
	//fmt.Println("Creating", path)  change to log
	err := os.MkdirAll(filepath.Dir(path), os.ModeDir|os.ModePerm)
	Check(err)
	
	// Clone if item is a file
	rc, err := f.Open()
	Check(err)
	if !f.FileInfo().IsDir() {	
		// Use os.Create() since Zip don't store file permissions.
		fileCopy, err := os.Create(path)
		Check(err)
		_, err = io.Copy(fileCopy, rc)
		fileCopy.Close()
		Check(err)
	}
	rc.Close()
}

func Extract(zip_path, dest string) {
	r, err := zip.OpenReader(zip_path)
	Check(err)
	defer r.Close()
	for _, f := range r.File {
		cloneZipItem(f, dest)
	}
}
