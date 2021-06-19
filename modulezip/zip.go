package modulezip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var logger, _ = zap.NewDevelopment()

func ZipModule(module string, version string) (io.Reader, error) {
	logger.Info("Creating a zip file with module files... ", zap.String("module", module), zap.String("version", version))
	buf := new(bytes.Buffer)
	// Create a new zip archive.
	w := zip.NewWriter(buf)
	// Add some files to the archive.

	path.Join(os.TempDir(), module, version)

	err := addFiles(w, path.Join(os.TempDir(), module, version), fmt.Sprintf("%s@%s", module, version))
	if err != nil {
		logger.Error("Error while creating zip", zap.Error(err))
		return nil, err
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		logger.Error("Error closing the zip ! ", zap.Error(err))
		return nil, err
	}

	return buf, nil
}

func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		pathFile := path.Join(basePath, file.Name())
		pathZip := path.Join(baseInZip, file.Name())
		if strings.HasPrefix(file.Name(), ".git") {
			logger.Debug("Skipping .git file/folder !")
			continue
		}

		logger.Debug(fmt.Sprintf("Adding file: %s -> %s", pathFile, pathZip))

		if !file.IsDir() {
			dat, err := ioutil.ReadFile(pathFile)
			if err != nil {
				return err
			}

			// Add some files to the archive.
			f, err := w.Create(pathZip)
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Recurse
			newBase := pathFile + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			err := addFiles(w, newBase, pathZip+"/")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
