package secret

import (
	"crypto/hmac"
	"crypto/sha256"
)

const algoHmacSha256 = "hmac-sha256"

type HmacSha256 struct{}

func (h *HmacSha256) Sign(msg string, secret string) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	if _, err := mac.Write([]byte(msg)); err != nil {
		return nil, err
	}

	return mac.Sum(nil), nil
}

func (h *HmacSha256) Name() string {
	return algoHmacSha256
}
