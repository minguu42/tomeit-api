package tomeit

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type user struct {
	id        int64
	digestUID string
	restCount int
	createdAt time.Time
}

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
