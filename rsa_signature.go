package helpers

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func parsePublicKey(pemBytes []byte) (Unsigner, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("ssh: no key found")
	}

	var rawKey interface{}

	switch block.Type {
	case "PUBLIC KEY":
		rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		rawKey = rsa
	default:
		return nil, errors.New("ssh: unsupported key type")
	}

	return newUnsignerFromKey(rawKey)
}

type Unsigner interface {
	Unsign(message []byte, signature []byte) error
}

func newUnsignerFromKey(k interface{}) (Unsigner, error) {
	var sshKey Unsigner

	switch t := k.(type) {
	case *rsa.PublicKey:
		sshKey = &rsaPublicKey{t}
	default:
		return nil, errors.New("ssh: unsupported key type")
	}

	return sshKey, nil
}

type rsaPublicKey struct {
	*rsa.PublicKey
}

func (r *rsaPublicKey) Unsign(message []byte, signature []byte) error {
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)

	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, signature)
}
