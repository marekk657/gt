package services

type ArchiveService interface {
	//
	CreateArchive(filePaths []string, destinationPath string) error

	// Extract extracts archive files into temp location.
	// Cleanup should be handled by caller
	Extract(archivePath string) ([]string, error)
}
