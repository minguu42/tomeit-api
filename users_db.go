package tomeit

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func (db *DB) createUser(digestUID string) (*User, error) {
	createdAt := time.Now()

	user := User{
		DigestUID: digestUID,
		RestCount: 4,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}

	r := db.Create(&user)
	if r.Error != nil {
		return nil, fmt.Errorf("db.Create failed: %w", r.Error)
	}

	return &user, nil
}

func (db *DB) getUserByDigestUID(digestUID string) (*User, error) {
	var user User

	r := db.Where("digest_uid", digestUID).First(&user)
	if errors.Is(r.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("record not found")
	}

	return &user, nil
}

func (db *DB) decrementRestCount(user *user) error {
	restCount := user.restCount

	if restCount == 1 {
		restCount = 4
	} else {
		restCount -= 1
	}

	const q = `UPDATE users SET rest_count = ? WHERE id = ?`

	if _, err := db.Exec(q, restCount, user.id); err != nil {
		return fmt.Errorf("db.Exec failed: %w", err)
	}

	return nil
}
