package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/ogreks/meeseeks-box/pkg/crypt"
	"github.com/stretchr/testify/assert"
)

func Test_RSACrypt(t *testing.T) {
	testCase := []struct {
		name    string
		value   string
		encode  crypt.Encode
		hash    crypt.Hash
		before  func(t *testing.T) *RSASecret
		wantVal bool
	}{
		{
			name: "test rsa scrypt encode/decode",
			before: func(t *testing.T) *RSASecret {
				privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
				assert.NoError(t, err)

				x509PublicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
				assert.NoError(t, err)

				return &RSASecret{
					PrivateKeyType:     crypt.PKCS1,
					PrivateKey:         base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey)),
					PrivateKeyDataType: crypt.Base64,
					PublicKey:          base64.StdEncoding.EncodeToString(x509PublicKey),
					PublicKeyDataType:  crypt.Base64,
				}
			},
			encode:  crypt.Base64,
			hash:    crypt.SHA256,
			value:   "hello world",
			wantVal: true,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			secret := tc.before(t)
			crypts := NewRSACrypt(secret)
			dst, err := crypts.Encrypt([]byte(tc.value), tc.encode)
			assert.NoError(t, err)
			t.Logf("\nencrypted data: %s", dst)

			text, err := crypts.Decrypt(dst, tc.encode)
			assert.NoError(t, err)
			t.Logf("\ndecrypted data: %s", text)

			sign, err := crypts.Sign(tc.value, tc.hash, tc.encode)
			assert.NoError(t, err)
			t.Logf("\nsign data: %s", sign)

			result, err := crypts.VerifySign(tc.value, tc.hash, sign, tc.encode)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantVal, result)
		})
	}
}
