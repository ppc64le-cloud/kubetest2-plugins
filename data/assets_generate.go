// +build ignore

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shurcooL/vfsgen"
)

var inputFS http.FileSystem = modTimeFS{
	fs: http.Dir("data"),
}

func main() {

	err := vfsgen.Generate(inputFS, vfsgen.Options{
		PackageName:  "data",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}

type modTimeFS struct {
	fs http.FileSystem
}

func (fs modTimeFS) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return modTimeFile{f}, nil
}

type modTimeFile struct {
	http.File
}

func (f modTimeFile) Stat() (os.FileInfo, error) {
	fi, err := f.File.Stat()
	if err != nil {
		return nil, err
	}
	return modTimeFileInfo{fi}, nil
}

type modTimeFileInfo struct {
	os.FileInfo
}

func (modTimeFileInfo) ModTime() time.Time {
	return time.Time{}
}
