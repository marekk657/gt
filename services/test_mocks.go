package services

type ArchiveServiceMock struct {
	CreateArchiveFunc func(filePaths []string, destinationPath string) error
	ExtractFunc       func(archivePath string) ([]string, error)
}

func (m ArchiveServiceMock) CreateArchive(filePaths []string, destinationPath string) error {
	if m.CreateArchiveFunc == nil {
		panic("CreateArchiveFunc is uninitialized!")
	}
	return m.CreateArchiveFunc(filePaths, destinationPath)
}

func (m ArchiveServiceMock) Extract(archivePath string) ([]string, error) {
	if m.ExtractFunc == nil {
		panic("ExtractFunc is uninitialized!")
	}
	return m.ExtractFunc(archivePath)
}
