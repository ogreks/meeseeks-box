//go:build local_ignore
// +build local_ignore

package utils

import (
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//go:test -timeout 300s -run ^Test_LocalDriver$ github.com/ogreks/meeseeks-box/internal/pkg/utils/
func Test_GenerateApp(t *testing.T) {
	key, err := GetAppKey()
	assert.NoError(t, err)
	secret, err := GetSecret(key)
	assert.NoError(t, err)

	t.Logf("session no: %s, app key: %s, app secret: %s\n", xid.New().String(), key, secret)
	t.Logf("%s", time.Now().Format(time.RFC850))
}
