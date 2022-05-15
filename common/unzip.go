package common

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractZipFile(source, destination string) ([]string, error) {
	r, err := zip.OpenReader(source)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = os.Stat(destination); os.IsNotExist(err) {
		if err = os.MkdirAll(destination, 0o777); err != nil {
			return nil, err
		}
	}

	var extractedFiles []string
	for _, f := range r.File {
		err = extractAndWriteFile(destination, f)
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, f.Name)
	}

	return extractedFiles, nil
}

func extractAndWriteFile(destination string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	pathM := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(pathM, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", pathM)
	}

	if f.FileInfo().IsDir() {
		err = os.MkdirAll(pathM, 0o777)
		if err != nil {
			return err
		}
	} else {
		err = os.MkdirAll(filepath.Dir(pathM), 0o777)
		if err != nil {
			return err
		}

		subFile, err := os.OpenFile(pathM, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			if err := subFile.Close(); err != nil {
				panic(err)
			}
		}()

		_, err = io.Copy(subFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
