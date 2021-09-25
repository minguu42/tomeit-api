package tomeit

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type User struct {
	ID        int
	DigestUID string
	RestCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func hash(token string) string {
	bytes := sha256.Sum256([]byte(token))
	digestToken := hex.EncodeToString(bytes[:])
	return digestToken
}
