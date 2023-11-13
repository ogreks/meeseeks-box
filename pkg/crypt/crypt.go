package crypt

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
)

// GetHash gets the crypto hash type & hashed data in different hash type
func GetHash(data []byte, hashType Hash) (h crypto.Hash, hashed []byte, err error) {
	f, h := GetHashFunc(hashType)
	result := f()
	if _, err = result.Write(data); err != nil {
		return
	}

	hashed = result.Sum(nil)
	return
}

// GetHashFunc gets the crypto hash func & type in different hash type
func GetHashFunc(hashType Hash) (f func() hash.Hash, h crypto.Hash) {
	switch hashType {
	case SHA1:
		f = sha1.New
		h = crypto.SHA1
	case SHA224:
		f = sha256.New224
		h = crypto.SHA224
	case SHA256:
		f = sha256.New
		h = crypto.SHA256
	case SHA384:
		f = sha512.New384
		h = crypto.SHA384
	case SHA512:
		f = sha512.New
		h = crypto.SHA512
	case SHA512_224:
		f = sha512.New512_224
		h = crypto.SHA512_224
	case SHA512_256:
		f = sha512.New512_256
		h = crypto.SHA512_256
	case MD5:
		f = md5.New
		h = crypto.MD5
	default:
		panic("unsupport hashType")
	}
	return
}

// DecodeString decodes string data to bytes in designed encoded type
func DecodeString(data string, encodedType Encode) ([]byte, error) {
	var (
		keyDecoded []byte
		err        error
	)

	switch encodedType {
	case String:
		keyDecoded = []byte(data)
	case Hex:
		keyDecoded, err = hex.DecodeString(data)
	case Base64:
		keyDecoded, err = base64.StdEncoding.DecodeString(data)
	default:
		return keyDecoded, fmt.Errorf("secretInfo PublicKeyDataType unsupport")
	}

	return keyDecoded, err
}

// ParsePrivateKey parses private key data to *rsa.PrivateKey
func ParsePrivateKey(data []byte, keyType Secret) (*rsa.PrivateKey, error) {
	switch keyType {
	case PKCS1:
		return x509.ParsePKCS1PrivateKey(data)
	case PKCS8:
		keyParsed, err := x509.ParsePKCS8PrivateKey(data)
		return keyParsed.(*rsa.PrivateKey), err
	default:
		return nil, fmt.Errorf("secretInfo PrivateKeyType unsupport")
	}
}

// EncodeToString encodes data to string with encode type
func EncodeToString(data []byte, encodeType Encode) (string, error) {
	switch encodeType {
	case Hex:
		return hex.EncodeToString(data), nil
	case Base64:
		return base64.StdEncoding.EncodeToString(data), nil
	case String:
		return string(data), nil
	default:
		return "", fmt.Errorf("secretInfo OutputType unsupport")
	}
}

// UnPaddingPKCS7 un-padding src data to original data , adapt to PKCS5 & PKCS7
func UnPaddingPKCS7(src []byte, blockSize int) []byte {
	n := len(src)
	if n == 0 {
		return src
	}
	paddingNum := int(src[n-1])
	if n < paddingNum || paddingNum > blockSize {
		return src
	}
	return src[:n-paddingNum]
}

// PKCS7Padding adds padding data using pkcs7 rules , adapt to PKCS5 &PKCS7
func PKCS7Padding(src []byte, blockSize int) []byte {
	paddingNum := blockSize - len(src)%blockSize
	padding := bytes.Repeat([]byte{byte(paddingNum)}, paddingNum)
	return append(src, padding...)
}
