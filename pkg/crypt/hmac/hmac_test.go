package hmac

import (
	"testing"

	"github.com/ogreks/meeseeks-box/pkg/crypt"
	"github.com/stretchr/testify/assert"
)

func Test_HMAC(t *testing.T) {
	testCase := []struct {
		name       string
		src        []byte
		hashType   crypt.Hash
		key        []byte
		hashed     string
		encodeType crypt.Encode
	}{
		{
			name:       "test md5 hmac hash",
			key:        []byte("D2LOfHWU7xlf8JbR"),
			src:        []byte("123456"),
			hashType:   crypt.MD5,
			hashed:     "96a0f2ed8bcedd2eac0efdd685b5814c",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha1 hmac hash",
			key:        []byte("D2LOfHWU7xlf8JbR"),
			src:        []byte("123456"),
			hashType:   crypt.SHA1,
			hashed:     "ea8b6afdb446a9bf06ef4fd4da61ddcd8ef1f426",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha224 hmac hash",
			key:        []byte("D2LOfHWU7xlf8JbR"),
			src:        []byte("123456"),
			hashType:   crypt.SHA224,
			hashed:     "c158976429d7a36e6d6a8287afa79fde76f196d45fa65dca6cd1b4b2",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha256 hmac hash",
			key:        []byte("D2LOfHWU7xlf8JbR"),
			src:        []byte("123456"),
			hashType:   crypt.SHA256,
			hashed:     "44adac09dbcab9f2e06ca7fcb706b32317705c2d18cf554bfa42f01cde6e703a",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha384 hmac hash",
			key:        []byte("D2LOfHWU7xlf8JbR"),
			src:        []byte("123456"),
			hashType:   crypt.SHA384,
			hashed:     "2a761d256b7d4fb97ee0d319de01769408e0f122740ce3b1834364bfe8d530c77ce097547699da1f792743fa9d129a87",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha512 hmac hash",
			key:        []byte("D2LOfHWU7xlf8JbR"),
			src:        []byte("123456"),
			hashType:   crypt.SHA512,
			hashed:     "3addd2d322e2a2c308f105061d115246f081fd50a2afc39aed79f8e2b5dabe769e6b05259d28b77ec9f4539e86182f319cb8a6b61b01511fb20f583cd61ff49c",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha512_224 hmac hash",
			key:        []byte("D2LOfHWU7xlf8JbR"),
			src:        []byte("123456"),
			hashType:   crypt.SHA512_256,
			hashed:     "163da5721b6460f5a0d7c5a0e899ac0c45f1b2bb3e0d146bb5685aa488b667d8",
			encodeType: crypt.Hex,
		},
		{
			name:       "test sha512_224 hmac hash",
			key:        []byte("D2LOfHWU7xlf8JbR"),
			src:        []byte("123456"),
			hashType:   crypt.SHA512_224,
			hashed:     "b4aacbc9183194363da6357082245face6258e8bcf7a8dc472ed97f2",
			encodeType: crypt.Hex,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			hm := NewHMAC(tc.hashType, tc.key)
			dst, err := hm.EncodeToString(tc.src, tc.encodeType)
			assert.NoError(t, err)

			assert.Equal(t, tc.hashed, dst)
		})
	}
}

func Test_HelperGenerateHAMCValue(t *testing.T) {
	var (
		src = []byte("123456")
		key = []byte("D2LOfHWU7xlf8JbR")
	)

	hmac := NewHMAC(crypt.SHA1, key)
	dst, err := hmac.EncodeToString(src, crypt.Hex)
	assert.NoError(t, err)

	t.Logf("generate value: %s", dst)
}
