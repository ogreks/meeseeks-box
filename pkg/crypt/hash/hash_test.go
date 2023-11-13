package hash

import (
	"testing"

	"github.com/ogreks/meeseeks-box/pkg/crypt"
	"github.com/stretchr/testify/assert"
)

func Test_Hash(t *testing.T) {
	testCase := []struct {
		name       string
		src        []byte
		hashType   crypt.Hash
		hashed     string
		encodeType crypt.Encode
	}{
		{
			name:       "test md5 hash",
			src:        []byte("123456"),
			hashType:   crypt.MD5,
			hashed:     "e10adc3949ba59abbe56e057f20f883e",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha1 hash",
			src:        []byte("123456"),
			hashType:   crypt.SHA1,
			hashed:     "7c4a8d09ca3762af61e59520943dc26494f8941b",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha224 hash",
			src:        []byte("123456"),
			hashType:   crypt.SHA224,
			hashed:     "f8cdb04495ded47615258f9dc6a3f4707fd2405434fefc3cbf4ef4e6",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha256 hash",
			src:        []byte("123456"),
			hashType:   crypt.SHA256,
			hashed:     "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha384 hash",
			src:        []byte("123456"),
			hashType:   crypt.SHA384,
			hashed:     "0a989ebc4a77b56a6e2bb7b19d995d185ce44090c13e2984b7ecc6d446d4b61ea9991b76a4c2f04b1b4d244841449454",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha512 hash",
			src:        []byte("123456"),
			hashType:   crypt.SHA512,
			hashed:     "ba3253876aed6bc22d4a6ff53d8406c6ad864195ed144ab5c87621b6c233b548baeae6956df346ec8c17f5ea10f35ee3cbc514797ed7ddd3145464e2a0bab413",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha512_224 hash",
			src:        []byte("123456"),
			hashType:   crypt.SHA512_256,
			hashed:     "184b5379d5b5a7ab42d3de1d0ca1fedc1f0ffb14a7673ebd026a6369745deb72",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha512_256 hash",
			src:        []byte("123456"),
			hashType:   crypt.SHA512_224,
			hashed:     "007ca663c61310fbee4c1680a5bbe70071825079b23f092713383296",
			encodeType: crypt.Hex,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHash(tc.hashType)
			dst, err := h.EncodeToString(tc.src, tc.encodeType)
			assert.NoError(t, err)

			assert.Equal(t, tc.hashed, dst)
		})
	}
}

func Test_HelperGenerateHashValue(t *testing.T) {
	var (
		src        = []byte("123456")
		hashType   = crypt.MD5
		encodeType = crypt.Hex
	)
	h := NewHash(hashType)
	dst, err := h.EncodeToString(src, encodeType)
	assert.NoError(t, err)
	t.Logf("generate hash value: %s", dst)
}
