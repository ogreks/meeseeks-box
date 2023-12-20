package token

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_DefaultStore(t *testing.T) {
	file := "./test.db.json"
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe
	logger := zap.New(zapcore.NewTee(),
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)
	store, err := NewDefaultStore[string](file, logger)
	assert.NoError(t, err)
	//defer os.Remove(file)

	testCases := []struct {
		name   string
		before func(t *testing.T, store *DefaultStore[string])
		after  func(t *testing.T, store *DefaultStore[string])
		token  string
		expire time.Duration
	}{
		{
			name:   "set token",
			before: func(t *testing.T, store *DefaultStore[string]) {},
			after: func(t *testing.T, store *DefaultStore[string]) {
				assert.Equal(t, true, store.Exists("test"))
				assert.NoError(t, store.Delete("test"))
			},
			token:  "test",
			expire: 3 * time.Second,
		},
		{
			name:   "delete token",
			before: func(t *testing.T, store *DefaultStore[string]) {},
			after: func(t *testing.T, store *DefaultStore[string]) {
				assert.NoError(t, store.Delete("test"))
			},
			token:  "test",
			expire: 3 * time.Second,
		},
		{
			name: "exists token",
			before: func(t *testing.T, store *DefaultStore[string]) {
				assert.NoError(t, store.Set("test", 1*time.Second))
			},
			after: func(t *testing.T, store *DefaultStore[string]) {
				assert.Equal(t, false, store.Exists("test"))
				assert.NoError(t, store.Delete("test"))
			},
			token:  "test",
			expire: -3 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t, store)
			assert.NoError(t, store.Set(tc.token, tc.expire))
			tc.after(t, store)
			time.Sleep(3 * time.Second)
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	t.Logf("%v", <-store.Shutdown(ctx))
	assert.NoError(t, os.Remove(file))
}
