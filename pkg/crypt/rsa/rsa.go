package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/ogreks/meeseeks-box/pkg/crypt"
)

type RSACrypt struct {
	Secret *RSASecret
}

type RSASecret struct {
	PublicKey          string
	PublicKeyDataType  crypt.Encode
	PrivateKey         string
	PrivateKeyDataType crypt.Encode
	PrivateKeyType     crypt.Secret
}

// NewRSACrypt init with the RSA secret info
func NewRSACrypt(secret *RSASecret) *RSACrypt {
	return &RSACrypt{
		Secret: secret,
	}
}

// Encrypt encrypts the given message with public key
// src the original data
// outputDataType the encode type of encrypted data ,such as Base64,HEX
func (r *RSACrypt) Encrypt(src []byte, outputDataType crypt.Encode) (dst string, err error) {
	secret := r.Secret
	if secret.PublicKey == "" {
		return "", fmt.Errorf("secret PublicKey can't be empty")
	}

	pubKeyDecoded, err := crypt.DecodeString(secret.PublicKey, secret.PublicKeyDataType)
	if err != nil {
		return "", err
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyDecoded)
	if err != nil {
		return "", err
	}

	var dataEncrypted []byte
	dataEncrypted, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), src)
	if err != nil {
		return "", err
	}

	return crypt.EncodeToString(dataEncrypted, outputDataType)
}

// Decrypt decrypts a plaintext using private key
// src the encrypted data with public key
// srcType the encode type of encrypted data ,such as Base64,HEX
func (r *RSACrypt) Decrypt(src string, srcType crypt.Encode) (dst string, err error) {
	secret := r.Secret
	if secret.PrivateKey == "" {
		return "", fmt.Errorf("secretInfo PrivateKey can't be empty")
	}
	privateKeyDecoded, err := crypt.DecodeString(secret.PrivateKey, secret.PrivateKeyDataType)
	if err != nil {
		return
	}
	prvKey, err := crypt.ParsePrivateKey(privateKeyDecoded, secret.PrivateKeyType)
	if err != nil {
		return
	}
	decodeData, err := crypt.DecodeString(src, srcType)
	if err != nil {
		return
	}
	var dataDecrypted []byte
	dataDecrypted, err = rsa.DecryptPKCS1v15(rand.Reader, prvKey, decodeData)
	if err != nil {
		return
	}
	return string(dataDecrypted), nil

}

// Sign calculates the signature of input data with the hash type & private key
// src the original unsigned data
// hashType the type of hash ,such as MD5,SHA1...
// outputDataType the encode type of sign data ,such as Base64,HEX
func (rc *RSACrypt) Sign(src string, hashType crypt.Hash, outputDataType crypt.Encode) (dst string, err error) {
	secret := rc.Secret
	if secret.PrivateKey == "" {
		return "", fmt.Errorf("sercet PrivateKey can't be empty")
	}
	privateKeyDecoded, err := crypt.DecodeString(secret.PrivateKey, secret.PrivateKeyDataType)
	if err != nil {
		return
	}
	prvKey, err := crypt.ParsePrivateKey(privateKeyDecoded, secret.PrivateKeyType)
	if err != nil {
		return
	}
	cryptoHash, hashed, err := crypt.GetHash([]byte(src), hashType)
	if err != nil {
		return
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, prvKey, cryptoHash, hashed)
	if err != nil {
		return
	}
	return crypt.EncodeToString(signature, outputDataType)
}

// VerifySign verifies input data whether match the sign data with the public key
// src the original unsigned data
// signedData the data signed with private key
// hashType the type of hash ,such as MD5,SHA1...
// signDataType the encode type of sign data ,such as Base64,HEX
func (rc *RSACrypt) VerifySign(src string, hashType crypt.Hash, signedData string, signDataType crypt.Encode) (bool, error) {
	secret := rc.Secret
	if secret.PublicKey == "" {
		return false, fmt.Errorf("secretInfo PublicKey can't be empty")
	}
	publicKeyDecoded, err := crypt.DecodeString(secret.PublicKey, secret.PublicKeyDataType)
	if err != nil {
		return false, err
	}
	pubKey, err := x509.ParsePKIXPublicKey(publicKeyDecoded)
	if err != nil {
		return false, err
	}
	cryptoHash, hashed, err := crypt.GetHash([]byte(src), hashType)
	if err != nil {
		return false, err
	}
	signDecoded, err := crypt.DecodeString(signedData, signDataType)
	if err != nil {
		return false, err
	}
	if err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), cryptoHash, hashed, signDecoded); err != nil {
		return false, err
	}
	return true, nil
}
