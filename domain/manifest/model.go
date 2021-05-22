package manifest

// Model defines Manifest structure
type Model struct {
	Files        []DataFile `json:"files"`
	SignatureUri string     `json:"signature_uri"`
}

// DataFile points to associated file in container
type DataFile struct {
	Uri           string `json:"uri"`
	HashAlgorithm string `json:"hash_algorithm"`
	Hash          string `json:"hash"`
}
