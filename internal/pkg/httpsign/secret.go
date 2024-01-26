package httpsign

import (
	"context"
	"github.com/ogreks/meeseeks-box/internal/pkg/httpsign/secret"
	"sync"
)

type KeyID string

type Secret struct {
	Key       string
	Algorithm secret.Crypto
}

type Secrets struct {
	sync.RWMutex

	container map[KeyID]*Secret
}

func NewSecrets() *Secrets {
	return &Secrets{
		container: make(map[KeyID]*Secret),
	}
}

type EncryptionKey interface {
	Secret(ctx context.Context, keyID KeyID) (*Secret, error)
}
