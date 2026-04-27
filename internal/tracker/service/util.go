package service

import (
	"crypto/rand"
	"encoding/hex"
)

func newID(prefix string) string {
	var raw [8]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return prefix + "_fallback"
	}
	return prefix + "_" + hex.EncodeToString(raw[:])
}
