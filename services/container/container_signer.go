package container

import (
	"fmt"
	"gt/services"
	"os"
	"strings"
)

type Signer struct {
	sigCreator     SignatureCreator
	archiveService services.ArchiveService
}

func NewSigner(sigCreator SignatureCreator, archiveService services.ArchiveService) Signer {
	return Signer{
		sigCreator:     sigCreator,
		archiveService: archiveService,
	}
}

func (s Signer) AddSignature(containerPath string) error {
	filePaths, err := s.archiveService.Extract(containerPath)
	if err != nil {
		return err
	}

	signatureFiles := s.filterFilePathsByPrefix(filePaths, signatureFileExtension)
	newManifestName := fmt.Sprintf(manifestFileNamePattern, len(signatureFiles)+1)
	dataFilePaths := s.filterFilePathsNotContaining(filePaths, metaInfPathZip)

	resp, err := s.sigCreator.NewSignature(dataFilePaths, newManifestName)
	if err != nil {
		return err
	}

	filePaths = append(filePaths, resp.ManifestFilePath, resp.SignatureFilePath)

	return s.archiveService.CreateArchive(filePaths, containerPath)
}

// RemoveSignature removes specified signature by id. If no such signature found, error is returned.
func (s Signer) RemoveSignature(containerPath string, signatureID int) error {
	filePaths, err := s.archiveService.Extract(containerPath)
	if err != nil {
		return err
	}

	defer os.RemoveAll(tmpFolderPath)

	fileName := fmt.Sprintf(manifestFileNamePattern, signatureID)
	filteredManifestFile := s.filterFilePathsByPrefix(filePaths, fileName)
	if len(filteredManifestFile) == 0 {
		return fmt.Errorf("signature with id '%v' not found", signatureID)
	}

	filteredManifestFilePath := filteredManifestFile[0]
	if err := s.removeSignatureFiles(filteredManifestFilePath); err != nil {
		return err
	}

	newFileNames := s.filterFilePathsNotContaining(filePaths, filteredManifestFilePath)
	return s.archiveService.CreateArchive(newFileNames, containerPath)
}

func (s Signer) removeSignatureFiles(manifestPath string) error {
	if err := os.Remove(manifestPath); err != nil {
		return err
	}

	signatureFileName := fmt.Sprintf(signatureFileNamePattern, manifestPath)
	if err := os.Remove(signatureFileName); err != nil {
		return err
	}
	return nil
}

func (s Signer) filterFilePathsNotContaining(paths []string, notContains string) []string {
	var newFilePaths []string
	for _, fn := range paths {
		if !strings.Contains(fn, notContains) {
			newFilePaths = append(newFilePaths, fn)
		}
	}
	return newFilePaths
}

func (s Signer) filterFilePathsByPrefix(paths []string, prefix string) []string {
	filtered := make([]string, 0, len(paths))
	for _, fn := range paths {
		if strings.Contains(fn, prefix) {
			filtered = append(filtered, fn)
		}
	}
	return filtered
}
