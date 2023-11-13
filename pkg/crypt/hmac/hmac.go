package hmac

import (
	"crypto/hmac"

	"github.com/ogreks/meeseeks-box/pkg/crypt"
)

type HMAC struct {
	HashType crypt.Hash
	Key      []byte
}

func NewHMAC(hashType crypt.Hash, key []byte) *HMAC {
	return &HMAC{
		HashType: hashType,
		Key:      key,
	}
}

// Get gets hashed bytes with defined hashType & key
func (h *HMAC) Get(src []byte) (dst []byte, err error) {
	f, _ := crypt.GetHashFunc(h.HashType)
	hmac := hmac.New(f, h.Key)
	hmac.Write(src)
	dst = hmac.Sum(nil)
	return
}

// EncodeToString gets hashed bytes with defined hashType & key and then encode to string
func (h *HMAC) EncodeToString(src []byte, encodeType ...crypt.Encode) (dst string, err error) {
	hashed, err := h.Get(src)
	if err != nil {
		return
	}

	if len(encodeType) == 0 {
		encodeType[0] = crypt.Hex
	}

	return crypt.EncodeToString(hashed, encodeType[0])
}

// GetHash gets hashed bytes with defined hashType
func GetHash(src []byte, hashType crypt.Hash) (dst []byte, err error) {
	_, dst, err = crypt.GetHash(src, hashType)
	return
}

// GetHashEncodeToString gets hashed bytes with defined hashType and then encode to string
func GetHashEncodeToString(src []byte, hashType crypt.Hash, encodeType crypt.Encode) (dst string, err error) {
	hashed, err := GetHash(src, hashType)
	if err != nil {
		return
	}
	return crypt.EncodeToString(hashed, encodeType)
}

// GetHMACHash gets hmac hashed bytes with defined hashType & key
func GetHMACHash(src []byte, hashType crypt.Hash, key []byte) (dst []byte, err error) {
	h := NewHMAC(hashType, key)
	return h.Get(src)
}

// GetHMACHashEncodeToString gets hmac hashed bytes with defined hashType & key then encode to string
func GetHMACHashEncodeToString(src []byte, hashType crypt.Hash, key []byte, encodeType crypt.Encode) (dst string, err error) {
	h := NewHMAC(hashType, key)
	return h.EncodeToString(src, encodeType)
}
