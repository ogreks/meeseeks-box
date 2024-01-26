package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"github.com/ogreks/meeseeks-box/configs"
	"sort"
	"strconv"
	"strings"
)

var (
	chars = []string{"a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s",
		"t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5",
		"6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I",
		"J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V",
		"W", "X", "Y", "Z"}
)

// GetAppKey get key
func GetAppKey() (string, error) {
	uidB, err := uuid.NewUUID()
	if err != nil {
		return "", errors.New("new app key error")
	}
	u := strings.ReplaceAll(uidB.String(), "-", "")

	appKey := ""

	for i := 0; i < 8; i++ {
		str := u[i*4 : i*4+4]
		x, _ := strconv.ParseInt(str, 16, 64)
		appKey += chars[x%0x3e]
	}

	return appKey, nil
}

// GetSecret get secret
func GetSecret(key string) (string, error) {
	sli := []string{key, configs.ProjectName}
	sort.Strings(sli)
	str := strings.Join(sli, "")

	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil)), nil
}
