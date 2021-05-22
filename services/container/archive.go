package container

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ZipArchiveService struct{}

func NewZipArchiveService() ZipArchiveService {
	return ZipArchiveService{}
}

// CreateArchive has implementation to create zip archives.
func (zas ZipArchiveService) CreateArchive(filePaths []string, destinationPath string) error {
	zipFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, filePath := range filePaths {
		if err := zas.addFileToZip(zipWriter, filePath); err != nil {
			return err
		}
	}

	return nil
}

// Extract extracts archive. It saves extracted files into tmp directory.
// tmp directory has to be cleaned up by API caller after work is done.
func (zas ZipArchiveService) Extract(archivePath string) ([]string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var fileNames []string
	for _, f := range r.File {
		fpath := filepath.Join(tmpFolderPath, f.Name)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		fileNames = append(fileNames, fpath)

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return nil, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return nil, err
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return nil, err
		}
	}
	return fileNames, nil
}

func (zas ZipArchiveService) addFileToZip(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// preserve folder structure
	if strings.Contains(filePath, "manifest") {
		header.Name = fmt.Sprintf("%s%s", metaInfPathZip, filepath.Base(filePath))
	}

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}
