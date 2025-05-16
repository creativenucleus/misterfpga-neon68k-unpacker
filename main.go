package main

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Usage: "unzip a collection of X68000 games into folder ready to upload to MisterFPGA",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "src-zip",
				Usage:    "The source zip file (see info on Neon68K's website for this)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "src-collection",
				Usage:    "The collection folder (which of the subfolders in the zip file to extract)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dest-type",
				Usage:    "file or ftp",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "dest-folder",
				Usage: "The destination folder for the unzipped files",
			},
			&cli.StringFlag{
				Name:  "dest-ip",
				Usage: "The destination address of the MiSTer",
			},
		},
		Action: func(cCtx *cli.Context) error {
			srcZip := cCtx.String("src-zip")
			srcCollection := cCtx.String("src-collection")
			destType := cCtx.String("dest-type")
			if srcZip == "" || srcCollection == "" || destType == "" {
				log.Fatal("src-zip, src-collection, and dest-type are required")
			}

			var writer IWriter
			if destType == "file" {
				destFolder := cCtx.String("dest-folder")
				if destFolder == "" {
					log.Fatal("dest-folder is required when dest-type is 'file'")
				}

				destFolder, err := filepath.Abs(destFolder)
				if err != nil {
					return err
				}

				writer = fileWriter(destFolder)
			} else if destType == "ftp" {
				destIP := cCtx.String("dest-ip")
				if destIP == "" {
					log.Fatal("dest-ip is required when dest-type is 'ftp'")
				}

				conn, err := openFTPConnection(destIP)
				if err != nil {
					return err
				}
				defer conn.Quit()

				writer = ftpWriter(*conn)
			} else {
				log.Fatal("dest-type must be either 'file' or 'ftp'")
			}

			err := unzipCollectionFromOuterZip(srcZip, srcCollection, writer)
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

func unzipCollectionFromOuterZip(srcZipFile string, collectionFolder string, writer IWriter) error {
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

			innerZipFileReader, err := innerZipFile.Open()
			if err != nil {
				return err
			}
			defer innerZipFileReader.Close()

			err = writer(correctedInnerPath, innerZipFileReader)
			if err != nil {
				return err
			}
		}

		log.Println("Successfully unzipped:", correctedPath)
	}

	return nil
}

func openFTPConnection(misterIP string) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(misterIP+":21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = conn.Login("root", "1")
	if err != nil {
		return nil, err
	}

	return conn, nil
}
