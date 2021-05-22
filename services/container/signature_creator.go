package container

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gt/domain/manifest"
	"gt/services"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/guardtime/goksi/hash"
)

type SignatureCreatorResponse struct {
	ManifestFilePath  string
	SignatureFilePath string
}

type SignatureCreator interface {
	NewSignature(filePaths []string, manifestName string) (SignatureCreatorResponse, error)
}

type signatureCreator struct {
	ksiSigner services.KSISigner
}

func NewSignatureCreator(ksiSigner services.KSISigner) SignatureCreator {
	return signatureCreator{
		ksiSigner: ksiSigner,
	}
}

func (sc signatureCreator) NewSignature(filePaths []string, manifestName string) (SignatureCreatorResponse, error) {
	if err := sc.createManifest(filePaths, manifestName); err != nil {
		return SignatureCreatorResponse{}, err
	}

	manifestPath := filepath.Join(fullMetaInfPath, manifestName)

	sigPath, err := sc.createSignature(manifestPath)
	if err != nil {
		return SignatureCreatorResponse{}, err
	}

	return SignatureCreatorResponse{
		ManifestFilePath:  manifestPath,
		SignatureFilePath: sigPath,
	}, nil
}

func (sc signatureCreator) createManifest(filePaths []string, manifestName string) error {
	manifestModel := manifest.Model{
		Files:        make([]manifest.DataFile, 0, len(filePaths)),
		SignatureUri: fmt.Sprintf("%s%s.sig", metaInfPathZip, manifestName),
	}

	for _, fp := range filePaths {
		hash, alg, err := sc.createHash(fp)
		if err != nil {
			return nil
		}

		dataFile := manifest.DataFile{
			Uri:           filepath.Base(fp),
			Hash:          hash,
			HashAlgorithm: alg,
		}

		manifestModel.Files = append(manifestModel.Files, dataFile)
	}

	fullPath := filepath.Join(fullMetaInfPath, manifestName)
	f, err := os.OpenFile(fullPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(manifestModel, "", " ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fullPath, b, 0777)
}

func (sc signatureCreator) createHash(filePath string) (string, string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", "", nil
	}
	defer f.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, f); err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), "SHA256", nil
}

func (sc signatureCreator) createSignature(manifestPath string) (string, error) {
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return "", err
	}

	hsr, err := hash.Default.New()
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(hsr, manifestFile); err != nil {
		return "", err
	}

	manifestHash, err := hsr.Imprint()
	if err != nil {
		return "", err
	}

	sig, err := sc.ksiSigner.Sign(manifestHash)
	if err != nil {
		return "", err
	}

	bin, err := sig.Serialize()
	if err != nil {
		return "", err
	}

	sigFileName := fmt.Sprintf(signatureFileFullPathPattern, fullMetaInfPath, filepath.Base(manifestPath))
	sigFile, err := os.Create(sigFileName)
	if err != nil {
		return "", err
	}
	defer sigFile.Close()

	if _, err := sigFile.Write(bin); err != nil {
		return "", err
	}

	return sigFileName, nil
}
