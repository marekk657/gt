package container_test

import (
	"errors"
	"gt/services"
	"gt/services/container"
	"testing"
)

func TestCreate_SignatureCreationFails(t *testing.T) {
	expectedErr := errors.New("signature creation failure")
	sigCreator := container.SignatureCreatorMock{
		NewSignatureFunc: func(filePaths []string, manifestName string) (container.SignatureCreatorResponse, error) {
			return container.SignatureCreatorResponse{}, expectedErr
		},
	}

	var archiveService services.ArchiveServiceMock

	creator := container.NewCreator(sigCreator, archiveService)

	// Act
	err := creator.Create([]string{"file1.txt"}, "container.zip")

	// Assert
	if err != expectedErr {
		t.Fatalf("expected error '%s' but received '%s'", expectedErr, err)
	}
}

func TestCreate_CreateArchiveFails(t *testing.T) {
	sigCreator := container.SignatureCreatorMock{
		NewSignatureFunc: func(filePaths []string, manifestName string) (container.SignatureCreatorResponse, error) {
			return container.SignatureCreatorResponse{
				ManifestFilePath:  "manifest1.json",
				SignatureFilePath: "signature1.json.sig",
			}, nil
		},
	}

	expectedErr := errors.New("archive creation failure")
	archiveService := services.ArchiveServiceMock{
		CreateArchiveFunc: func(filePaths []string, destinationPath string) error {
			if len(filePaths) != 3 {
				t.Error(filePaths)
				t.Errorf("invalid count of file paths! got=%v, want=%v", len(filePaths), 3)
			}

			return expectedErr
		},
	}

	creator := container.NewCreator(sigCreator, archiveService)

	// Act
	err := creator.Create([]string{"file1.txt"}, "container.zip")

	// Assert
	if err != expectedErr {
		t.Fatalf("expected error '%s' but received '%s'", expectedErr, err)
	}
}
