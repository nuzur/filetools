package filetools

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ZipRequest struct {
	// zip file identifier
	Identifier string

	// output path
	OutputPath string
}

func GenerateZip(ctx context.Context, req ZipRequest) error {
	// creating new zip file
	rootDir := CurrentPath()
	zipFile := path.Join(rootDir, req.OutputPath, req.Identifier)
	archive, err := os.Create(fmt.Sprintf("%s.zip", zipFile))
	if err != nil {
		return err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		parts := strings.Split(path, req.Identifier)
		if len(parts) < 2 {
			return errors.New("invalid path")
		}
		f, err := zipWriter.Create(parts[1])
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}

	finalOutput := path.Join(rootDir, req.OutputPath, req.Identifier)
	err = filepath.Walk(finalOutput, walker)
	if err != nil {
		return err
	}

	return nil
}
