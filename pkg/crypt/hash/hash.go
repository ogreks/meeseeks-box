package hash

import "github.com/ogreks/meeseeks-box/pkg/crypt"

type Hash struct {
	HashType crypt.Hash
}

func NewHash(hashType crypt.Hash) *Hash {
	return &Hash{
		HashType: hashType,
	}
}

//Get gets hashed bytes with defined hashType
func (h *Hash) Get(src []byte) (dst []byte, err error) {
	_, dst, err = crypt.GetHash(src, h.HashType)
	return
}

//EncodeToString gets hashed bytes with defined hashType and then encode to string
func (h *Hash) EncodeToString(src []byte, encodeType ...crypt.Encode) (dst string, err error) {
	hashed, err := h.Get(src)
	if err != nil {
		return
	}

	if len(encodeType) == 0 {
		encodeType[0] = crypt.Hex
	}

	return crypt.EncodeToString(hashed, encodeType[0])
}
