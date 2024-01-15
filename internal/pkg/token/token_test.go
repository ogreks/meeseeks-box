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
			wantVal: "123123121",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tk := NewDefaultToken[string, func() (jwt.SigningMethod, []byte)](
				WithStore[string, func() (jwt.SigningMethod, []byte)](store),
				WithExpire[string, func() (jwt.SigningMethod, []byte)](100*time.Second),
				WithFun[string, func() (jwt.SigningMethod, []byte)](func() (jwt.SigningMethod, []byte) {
					return jwt.SigningMethodHS512, []byte("1231231231")
				}),
				WithClaims[string, func() (jwt.SigningMethod, []byte)](&JWTClaim{}),
			)

			tks, err := tk.CreateToken(context.TODO(), &JWTClaim{
				Content: tc.token,
			})
			assert.NoError(t, err)

			claims, err := tk.Validate(tks)
			assert.Equal(t, err, tc.wantErr)

			assert.Equal(t, claims.(*JWTClaim).Content.(string), tc.wantVal)
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Logf("%v", <-store.Shutdown(ctx))
	assert.NoError(t, os.Remove(file))
}
