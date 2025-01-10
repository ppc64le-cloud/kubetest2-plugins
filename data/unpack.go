package data

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	//go:embed k8s-ansible powervs vpc config.tf
	dir embed.FS
)

// Unpack handles copying out the embedded files from the binary to the destination.
// Accepts extractPath, which is the directory to extract to, on host.
// resPath holds the resource to be copied over from the binary to the host.

func Unpack(extractPath, resPath string) error {
	file, err := dir.Open(resPath)

	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.IsDir() {
		if err := os.MkdirAll(extractPath, 0777); err != nil {
			return fmt.Errorf("cannot create directory - %v", err)
		}
		file.Close()

		files, err := dir.ReadDir(resPath)
		if err != nil {
			return fmt.Errorf("cannot read the provided directory - %v", err)
		}
		for _, file := range files {
			name := file.Name()
			if err := Unpack(filepath.Join(extractPath, name), filepath.Join(resPath, name)); err != nil {
				return err
			}
		}
		return nil
	}

	out, err := os.OpenFile(extractPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}
