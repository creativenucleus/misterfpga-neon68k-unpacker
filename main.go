package main

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
		Usage: "unzip a collection of X68000 games into folder ready to upload to MisterFPGA",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src-zip",
				Usage:    "The source zip file (see info on Neon68K's website for this)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "collection",
				Usage:    "The collection folder (which of the subfolders in the zip file to extract)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dest-folder",
				Usage:    "The destination folder for the unzipped files",
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

			err := unzipCollectionFromOuterZip(srcZip, collection, destFolder)
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

func unzipCollectionFromOuterZip(srcZipFile string, collectionFolder string, destFolder string) error {
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

		// Use only those that match the requested collection path
		if !strings.HasPrefix(correctedPath, collectionFolder) {
			continue
		}

		log.Println("Unzipping: ", correctedPath)

		outerZipFileReader, err := outerZipFile.Open()
		if err != nil {
			return err
		}
		defer outerZipFileReader.Close()

		outerZipFileBuffer, err := io.ReadAll(outerZipFileReader)
		if err != nil {
			return err
		}

		innerZipReader, err := zip.NewReader(bytes.NewReader(outerZipFileBuffer), int64(len(outerZipFileBuffer)))
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

			destinationFile, err := os.OpenFile(correctedInnerPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
			if err != nil {
				return err
			}

			innerZipFileReader, err := innerZipFile.Open()
			if err != nil {
				return err
			}
			if _, err := io.Copy(destinationFile, innerZipFileReader); err != nil {
				return err
			}
		}

		log.Println("Successfully unzipped:", correctedPath)
	}

	return nil
}
