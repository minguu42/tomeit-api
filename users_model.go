package tomeit

import (
	"crypto/sha256"
	"encoding/hex"
)

type user struct {
	id        int64
	digestUID string
	restCount int
}

func hash(token string) string {
	bytes := sha256.Sum256([]byte(token))
	digestToken := hex.EncodeToString(bytes[:])
	return digestToken
}
