package token

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type JWTClaim struct {
	jwt.StandardClaims

	Content any `json:"content"`
}

func Test_DefaultCreateToken(t *testing.T) {
	file := "./test.db.json"
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe
	logger := zap.New(zapcore.NewTee(),
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)
	store, err := NewDefaultStore[string](file, logger)
	assert.NoError(t, err)

	testCase := []struct {
		name    string
		token   string
		wantVal string
		wantErr error
	}{
		{
			name:    "create token",
			token:   "123123121",
			wantVal: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjb250ZW50IjpudWxsfQ.hUOIv_i_PmoIJ4zI6qUzPkp5HzJa5yEEQl8WT28NMn45x6yHEFEQH5S5IWA37KMF8pQbqiyuAbH76MtcUcECbw",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tk := NewDefaultToken[string, func() (jwt.SigningMethod, []byte, jwt.Claims)](
				WithStore[string, func() (jwt.SigningMethod, []byte, jwt.Claims)](store),
				WithExpire[string, func() (jwt.SigningMethod, []byte, jwt.Claims)](100),
				WithFun[string, func() (jwt.SigningMethod, []byte, jwt.Claims)](func() (jwt.SigningMethod, []byte, jwt.Claims) {
					return jwt.SigningMethodHS512, []byte("1231231231"), &JWTClaim{}
				}),
			)

			val, err := tk.CreateToken(context.TODO())
			assert.Equal(t, val, tc.wantVal)
			assert.Equal(t, err, tc.wantErr)
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Logf("%v", <-store.Shutdown(ctx))
	assert.NoError(t, os.Remove(file))
}
