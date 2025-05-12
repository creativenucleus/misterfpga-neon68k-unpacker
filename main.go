package main

// Beware - this is quickly built and could be destructive to file system!

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "unzip",
		Usage: "unzip a load of zip files into a folder",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src-zip",
				Usage:    "src zip file",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "collection",
				Usage:    "collection folder",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dest-folder",
				Usage:    "dest folder",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			srcZip := cCtx.String("src-zip")
			collection := cCtx.String("collection")
			destFolder := cCtx.String("dest-folder")
			if srcZip == "" || collection == "" || destFolder == "" {
				log.Fatal("src-zip, collection, and dest-folder are required")
			}

			err := unzipFromOuterZip(srcZip, collection, destFolder)
			if err != nil {
				log.Fatal(err)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func unzipFromOuterZip(srcZipFile string, collectionFolder string, destFolder string) error {
	destFolder, err := filepath.Abs(destFolder)
	if err != nil {
		return err
	}

	outerZipReader, err := zip.OpenReader(srcZipFile)
	if err != nil {
		return err
	}
	defer outerZipReader.Close()

	for _, outerZipFile := range outerZipReader.File {
		correctedPath := filepath.FromSlash(outerZipFile.Name)

		if strings.HasPrefix(correctedPath, collectionFolder) {
			log.Println("Unzipping: ", correctedPath)

			outerZipReader, err := outerZipFile.Open()
			if err != nil {
				return err
			}
			defer outerZipReader.Close()

			buf, err := io.ReadAll(outerZipReader)
			if err != nil {
				return err
			}

			innerZipReader, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
			if err != nil {
				return err
			}

			for _, innerZipFile := range innerZipReader.File {
				correctedInnerPath := filepath.FromSlash(innerZipFile.Name)
				correctedInnerPath = filepath.Join(destFolder, correctedInnerPath)

				err := os.MkdirAll(filepath.Dir(correctedInnerPath), os.ModePerm)
				if err != nil {
					return err
				}

				dstFile, err := os.OpenFile(correctedInnerPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
				if err != nil {
					return err
				}

				fileInArchive, err := innerZipFile.Open()
				if err != nil {
					return err
				}
				if _, err := io.Copy(dstFile, fileInArchive); err != nil {
					return err
				}
			}

			log.Println("Successfully unzipped:", correctedPath)
		}
	}

	return nil
}
