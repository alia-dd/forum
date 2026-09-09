package utils

import (
	"crypto/rand"
	"fmt"
	"time"

	customerrors "gitea.kood.tech/jyrkikarhunen/forum/errors"
)

// uuid resourse https://www.rfc-editor.org/rfc/rfc9562.html?utm_source=chatgpt.com#name-uuid-version-7

func NewUuid() (string, error) {
	timestamp := time.Now().UnixMilli()

	b := make([]byte, 16)
	// generate 16 random bytes
	if _, randomizeErr := rand.Read(b); randomizeErr != nil {
		return "", customerrors.ErrInternalError
	}

	// uuid v7 includes timestamp at the most significant 48 bits
	// replace those bits with the timestamp
	b[0] = byte(timestamp >> 40)
	b[1] = byte(timestamp >> 32)
	b[2] = byte(timestamp >> 24)
	b[3] = byte(timestamp >> 16)
	b[4] = byte(timestamp >> 8)
	b[5] = byte(timestamp)

	b[6] = (b[6] & 0x0f) | 0x70 // Version 7 random
	b[8] = (b[8] & 0x3f) | 0x80 // Variant is 10

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}
