package main

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jlaffaye/ftp"
)

type IWriter func(destPath string, readCloser io.ReadCloser) error

// fileWriter returns a function that satisfies the IWriter interface and works with the local filesystem
func fileWriter(destFolder string) IWriter {
	return func(innerPath string, readCloser io.ReadCloser) error {
		fullPath := filepath.Join(destFolder, innerPath)
		err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
		if err != nil {
			return err
		}

		destinationFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}

		if _, err := io.Copy(destinationFile, readCloser); err != nil {
			return err
		}

		return nil
	}
}

// ftpMakeDirRecursive creates a directory and all its parents on the FTP server
func ftpMakeDirRecursive(conn ftp.ServerConn, dir string) error {
	// Split the FTP directory into its components
	parts := strings.Split(dir, "/")

	subDir := "/" + parts[0]
	for _, part := range parts[1:] {
		subDir = path.Join(subDir, part)

		err := conn.MakeDir(subDir)
		if err != nil && err.Error()[0:4] != "550 " { // File already exists error is fine
			return err
		}
	}

	return nil
}

// ftpWriter returns a function that satisfies the IWriter interface and works with FTP
func ftpWriter(conn ftp.ServerConn) IWriter {
	return func(innerPath string, readCloser io.ReadCloser) error {
		fullPath := path.Clean(path.Join("/media/fat/", filepath.ToSlash(innerPath)))
		dir, filename := filepath.Split(fullPath)

		err := ftpMakeDirRecursive(conn, dir)
		if err != nil {
			return err
		}

		err = conn.ChangeDir(dir)
		if err != nil {
			return err
		}

		err = conn.Stor(filename, readCloser)
		if err != nil {
			return err
		}

		return nil
	}
}
