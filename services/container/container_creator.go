package container

import (
	"gt/services"
	"os"
)

type Creator struct {
	sigCreator     SignatureCreator
	archiveService services.ArchiveService
}

func NewCreator(sigCreator SignatureCreator, archiveService services.ArchiveService) Creator {
	return Creator{
		sigCreator:     sigCreator,
		archiveService: archiveService,
	}
}

// Create creates new container to given "containerFullPath". And signs it's content.
// FilePaths slice contains all files that are added to container.
func (c Creator) Create(filePaths []string, containerFullPath string) error {
	os.MkdirAll(fullMetaInfPath, 0777)

	defer os.RemoveAll(tmpFolderPath)

	sigResponse, err := c.sigCreator.NewSignature(filePaths, initialManifestName)
	if err != nil {
		return err
	}

	filePaths = append(filePaths, sigResponse.ManifestFilePath, sigResponse.SignatureFilePath)

	return c.archiveService.CreateArchive(filePaths, containerFullPath)
}
