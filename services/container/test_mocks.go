package container

type SignatureCreatorMock struct {
	NewSignatureFunc func(filePaths []string, manifestName string) (SignatureCreatorResponse, error)
}

func (m SignatureCreatorMock) NewSignature(filePaths []string, manifestName string) (SignatureCreatorResponse, error) {
	if m.NewSignatureFunc == nil {
		panic("NewSignatureFunc is uninitialized!")
	}
	return m.NewSignatureFunc(filePaths, manifestName)
}
