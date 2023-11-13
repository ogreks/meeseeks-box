package aes

import (
	"crypto/aes"

	"github.com/ogreks/meeseeks-box/pkg/crypt"
)

type AESCrypt struct {
	crypt.CipherCrypt
}

func NewAESCrypt(key []byte) (*AESCrypt, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &AESCrypt{
		crypt.CipherCrypt{
			Block: block,
		},
	}, nil
}
