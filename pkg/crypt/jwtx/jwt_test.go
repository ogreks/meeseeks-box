package jwtx

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func Test_JWT(t *testing.T) {
	testCase := []struct {
		name          string
		key           KeyFunc
		val           string
		secret        Secret
		signingMethod SigningMethod
	}{
		{
			name: "test jwt hs",
			key: func(t *Token) (interface{}, error) {
				return []byte("a40065df93ff3e35febb16e66f0dda29"), nil
			},
			val: "hello world",
			secret: func() []byte {
				return []byte("a40065df93ff3e35febb16e66f0dda29")
			},
			signingMethod: jwt.SigningMethodHS256,
		},
		// TODO 不支持 rsa 签名
		// {
		// 	name: "test jwt rsa",
		// 	key: func(t *Token) (interface{}, error) {
		// 		return []byte(`MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtVXbyqanOBbQZbQWrPpr6TAp1gKuvKVZVvgP8prPu5DfW2DxRZcrVUMiKuw1mmO8gRv5fn0hW2Q1gSgwt8rk9sWB4yWJjJVYldUHTKQ2pLbn+hA5Smq4zsZRiXEi6/F8j/RLKh6RgwVzW5mzivizvh88wXZ2NVO0ol5E4aRAx4jXrSRRUQySuh/N9H1bF4WZFqKlTKOv39KLe9IoyIba/FSOgThmC203si0D1/jTof/SR7pu4+IIXJOPatdfqnXOLZJzorlSkbGc5iM0oiD8aqbD8kD5xZM/gxRIVdNuy5bMtqJibAQri4CSCtWQ78NZRJYggJgJ2LrPLRy2RyRb6wIDAQAB`), nil
		// 	},
		// 	val: "hello world",
		// 	secret: func() []byte {
		// 		return []byte(`MIIEowIBAAKCAQEAtVXbyqanOBbQZbQWrPpr6TAp1gKuvKVZVvgP8prPu5DfW2DxRZcrVUMiKuw1mmO8gRv5fn0hW2Q1gSgwt8rk9sWB4yWJjJVYldUHTKQ2pLbn+hA5Smq4zsZRiXEi6/F8j/RLKh6RgwVzW5mzivizvh88wXZ2NVO0ol5E4aRAx4jXrSRRUQySuh/N9H1bF4WZFqKlTKOv39KLe9IoyIba/FSOgThmC203si0D1/jTof/SR7pu4+IIXJOPatdfqnXOLZJzorlSkbGc5iM0oiD8aqbD8kD5xZM/gxRIVdNuy5bMtqJibAQri4CSCtWQ78NZRJYggJgJ2LrPLRy2RyRb6wIDAQABAoIBACmfapLqi+mI+w7NNoM/+/MLlh7EgN1WStp2mBqExHf2Of6ckuT5XP32KeqWS3uDtyofhLYu/LAgoVNjKUOWs9Wc4kKERD1brni97C4AZ3aJcVfpz2ywcHkt9ltI984WoRPd7D7fw2DCqIb3EcotafyS6PwzN9DnxMecQc1VSXVDJdDSMKdCQmUYSTfeiJx3L/CinaB2NuErJrTKDBdG5//+A/16LZWtHSz2Hjzx4tNOS6s9gqWnwqiSf/zdLi5FbKp8fG4zpfbf6ICMdv9uzINU5Aw4Hb4VV6xSaiB6cyL8E55wBbtlBjWRuJT3jV3iS/sBoXgbo25BAkxB4zMcg2kCgYEAzYQVbrF1YiF+zm4PCUqmkqzHGdgIK+BaVfkIu3oRcymj7WickES2aqYkUJxoiOWZ6hrVBWi7+tbHEKQbchOb0hRkHl01WkX9QDBJOgLRzM4rH7lo+HHnewCaWBCBHt+Pgsx2r3eIvdGoMW45BujLCfMkxPsg2eE/pWl70m/MJY8CgYEA4eEsrOyhzECz82ODs+zm+4Zr0WFys2/wsqCbObhCQxmjImgYrVx5/6PPvOXTT/CsyW6VBX+IRlKbOAs2EI2my8k+0yh7ThLvX65DFtKazB9cUnddW6Mo+OnXaP9Z
		// 		6KIQrB29sUwTy+SPsQJxt2xrVgSdGjYO7Tlzsy3VfAg5jeUCgYBjpUc/BIwFqHfzXymrN6bMNznSLgzIOV+Tj8vMGsObMvVohvBigu5vh17UNwH4XlriR2BB3yZF2R1r1CX6icdjdL+WeVsFCipglQjsN7HBu7TtDNj4nUG/QeuUB2yTq5HJuZlSOoLyhlUmomrEDttjV6DcYWbsPWq+qQaAYfR5wQKBgQChTBmJ9oRnhcC6ydJPlpku/cNaRjsRFZuNAxE1e3Wd0t3igPE2QrY/ret3Waq3CAdq5BN4VKSsiuqab68QzQZRuYiqYtsCWpUi/x6bWpL9tltH7EL3YCCu9tVC/i1m6Ov87FP8GnZ8f994KGWp9LsFNtA02mt4TTFovw8WvgzTXQKBgEVHG7E+phgEzpG4sm306qhzs2FvLfF8AOQS6yAwqt5XmXz/YK8wKx40Gng9KMHl39iisXzY9Jz1E1uEwfVx5rWmcp76RvgYGwRz+G1B6jE4F2bPlsDIqplsGktc7pZkSgzS4TfM6aTZoHXqIJASpKCRuXkRVR23Yc/fNr6kkBjQ`)
		// 	},
		// 	signingMethod: jwt.SigningMethodRS256,
		// },
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			j := NewJWT(tc.signingMethod)
			j.BuildClaims(tc.val, "mark", "issuer", 0)
			token, err := j.GetToken()
			assert.NoError(t, err)
			t.Logf("\ntoken: %s", token)

			claims, err := j.ParseToken(token, tc.key)
			assert.NoError(t, err)

			assert.Equal(t, claims.Content.(string), tc.val)

			t.Logf("\nclaims: %+v", claims)
		})
	}
}

func Test_HelpGenerateKey(t *testing.T) {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	assert.NoError(t, err)

	t.Log(hex.EncodeToString(key))
}

func Test_HelpGenerateRsaKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	x509PublicKey, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	assert.NoError(t, err)

	t.Logf("\nprivate key: \n%s\n", base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PrivateKey(privateKey)))
	t.Logf("\npublic key: \n%s\n", base64.StdEncoding.EncodeToString(x509PublicKey))
}
