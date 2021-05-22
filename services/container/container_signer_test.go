package container_test

import (
	"errors"
	"gt/services"
	"gt/services/container"
	"os"
	"testing"
)

func TestAddSignature_ExtractFails(t *testing.T) {
	expectedErr := errors.New("failed to extract")
	archiveService := services.ArchiveServiceMock{
		ExtractFunc: func(archivePath string) ([]string, error) {
			return nil, expectedErr
		},
	}

	var sigCreator container.SignatureCreatorMock

	signer := container.NewSigner(sigCreator, archiveService)

	// Act
	err := signer.AddSignature("container.zip")

	// Assert
	if err != expectedErr {
		t.Fatalf("expected error '%s' but received '%s'", expectedErr, err)
	}
}

func TestAddSignature_CreatingNewSignatureFails(t *testing.T) {
	archiveService := services.ArchiveServiceMock{
		ExtractFunc: func(archivePath string) ([]string, error) {
			return []string{"text1.txt", "META-INF/manifest1.json", "META-INF/manifest1.json.sig"}, nil
		},
	}

	expectedErr := errors.New("failed to sign")
	sigCreator := container.SignatureCreatorMock{
		NewSignatureFunc: func(filePaths []string, manifestName string) (container.SignatureCreatorResponse, error) {
			return container.SignatureCreatorResponse{}, expectedErr
		},
	}

	signer := container.NewSigner(sigCreator, archiveService)

	// Act
	err := signer.AddSignature("container.zip")

	// Assert
	if err != expectedErr {
		t.Fatalf("expected error '%s' but received '%s'", expectedErr, err)
	}
}

func TestAddSignature_CreatingArchiveFails(t *testing.T) {
	expectedErr := errors.New("failed to create archive")

	archiveService := services.ArchiveServiceMock{
		ExtractFunc: func(archivePath string) ([]string, error) {
			return []string{"text1.txt", "META-INF/manifest1.json", "META-INF/manifest1.json.sig"}, nil
		},
		CreateArchiveFunc: func(filePaths []string, destinationPath string) error {
			if destinationPath == "" {
				return errors.New("destination is empty")
			}

			if len(filePaths) != 5 {
				t.Error(filePaths)
				return errors.New("invalid count of file paths in creatArchive mock")
			}

			return expectedErr
		},
	}

	sigCreator := container.SignatureCreatorMock{
		NewSignatureFunc: func(filePaths []string, manifestName string) (container.SignatureCreatorResponse, error) {
			if len(filePaths) != 1 {
				t.Error(filePaths)
				return container.SignatureCreatorResponse{}, errors.New("invalid count of filepaths in signature creator mock")
			}

			if manifestName == "" {
				return container.SignatureCreatorResponse{}, errors.New("manifest name is empty")
			}

			return container.SignatureCreatorResponse{
				ManifestFilePath:  "META-INF/manifest2.json",
				SignatureFilePath: "META-INF/manifest2.json.sig",
			}, nil
		},
	}

	signer := container.NewSigner(sigCreator, archiveService)

	// Act
	err := signer.AddSignature("container.zip")

	// Assert
	if err != expectedErr {
		t.Fatalf("expected error '%s' but received '%s'", expectedErr, err)
	}
}

func TestRemoveSignature_ExtractFails(t *testing.T) {
	expectedErr := errors.New("failed to extract")
	archiveService := services.ArchiveServiceMock{
		ExtractFunc: func(archivePath string) ([]string, error) {
			return nil, expectedErr
		},
	}

	var sigCreator container.SignatureCreatorMock

	signer := container.NewSigner(sigCreator, archiveService)

	// Act
	err := signer.RemoveSignature("container.zip", 1)

	// Assert
	if err != expectedErr {
		t.Fatalf("expected error '%s' but received '%s'", expectedErr, err)
	}
}

func TestRemoveSignature_SignatureWithGivenIDNotFound(t *testing.T) {
	archiveService := services.ArchiveServiceMock{
		ExtractFunc: func(archivePath string) ([]string, error) {
			return []string{"text1.txt", "META-INF/manifest1.json", "META-INF/manifest1.json.sig"}, nil
		},
	}

	expectedErrStr := "signature with id '2' not found"
	var sigCreator container.SignatureCreatorMock

	signer := container.NewSigner(sigCreator, archiveService)

	// Act
	err := signer.RemoveSignature("container.zip", 2)

	// Assert
	if err.Error() != expectedErrStr {
		t.Fatalf("expected error '%s' but received '%s'", expectedErrStr, err)
	}
}

func TestRemoveSignature_CreateArchiveFails(t *testing.T) {
	testSetup(t)
	defer cleanup(t)

	expectedErr := errors.New("create archive error")
	archiveService := services.ArchiveServiceMock{
		ExtractFunc: func(archivePath string) ([]string, error) {
			return []string{"text1.txt", "META-INF/manifest1.json", "META-INF/manifest1.json.sig"}, nil
		},
		CreateArchiveFunc: func(filePaths []string, destinationPath string) error {
			if len(filePaths) != 1 {
				t.Error(filePaths)
				return errors.New("invalid count of filepaths in creatArchive mock")
			}

			return expectedErr
		},
	}

	var sigCreator container.SignatureCreatorMock

	signer := container.NewSigner(sigCreator, archiveService)

	// Act
	err := signer.RemoveSignature("container.zip", 1)

	// Assert
	if err != expectedErr {
		t.Fatalf("expected error '%s' but received '%s'", expectedErr, err)
	}
}

func testSetup(t *testing.T) {
	os.MkdirAll("META-INF", 0777)

	f, err := os.Create("META-INF/manifest1.json")
	if err != nil {
		t.Fatal("failed to setup test:", err)
	}
	defer f.Close()

	f2, err := os.Create("META-INF/manifest1.json.sig")
	if err != nil {
		t.Fatal("failed to setup test:", err)
	}
	defer f2.Close()
}

func cleanup(t *testing.T) {
	os.RemoveAll("META-INF")
}
