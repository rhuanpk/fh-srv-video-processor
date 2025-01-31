package handler

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"extractor/internal/config"
	"extractor/internal/resouce/zipper"
)

type service struct{}

func NewService() zipper.Service {
	return &service{}
}

func addFile(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}

	return nil
}

func (s service) Create(zipsPaths, sourcesDirs []string) error {
	for index, zipPath := range zipsPaths {
		zipFile, err := os.Create(zipPath)
		if err != nil {
			return err
		}
		defer zipFile.Close()

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		frames, err := filepath.Glob(filepath.Join(sourcesDirs[index], config.FrameName))
		if err != nil {
			return err
		}

		for _, frame := range frames {
			if err := addFile(zipWriter, frame); err != nil {
				return err
			}
		}
	}

	return nil
}
