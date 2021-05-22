package services

import (
	"github.com/guardtime/goksi/hash"
	"github.com/guardtime/goksi/service"
	"github.com/guardtime/goksi/signature"
)

// KSISigner is helper interface to wrap guardtime Signer struct.
type KSISigner interface {
	Sign(hash hash.Imprint, opt ...service.SignOption) (*signature.Signature, error)
}

// NewKSISigner creates Signer service, uses guardtime API underneath.
func NewKSISigner(endpoint, username, pswd string) (KSISigner, error) {
	return service.NewSigner(service.OptEndpoint(endpoint, username, pswd))
}
