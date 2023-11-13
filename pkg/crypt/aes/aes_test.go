package aes

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/ogreks/meeseeks-box/pkg/crypt"
	"github.com/stretchr/testify/assert"
)

func Test_AESCrypt(t *testing.T) {
	testCase := []struct {
		name  string
		key   string
		iv    string
		value string
	}{
		{
			name:  "test aes scrypt encode/decode",
			key:   "IctGGZUE0jlnuINCgg4T7A==",
			iv:    "6Wstcdw+KjMGRbL7pJroBw==",
			value: "hello world",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			key, _ := base64.StdEncoding.DecodeString(tc.key)
			iv, _ := base64.StdEncoding.DecodeString(tc.iv)
			value := []byte(tc.value)

			block, err := NewAESCrypt(key)
			assert.NoError(t, err)

			data, err := block.Encrypt(value, crypt.CBC, iv)
			assert.NoError(t, err)

			data, err = block.Decrypt(data, crypt.CBC, iv)
			assert.NoError(t, err)

			assert.Equal(t, value, data)
		})
	}
}

func Test_HelpGenerateIv(t *testing.T) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	assert.NoError(t, err)

	t.Log(hex.EncodeToString(key))
}
