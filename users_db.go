package tomeit

import (
	"fmt"
	"time"
)

type userDBInterface interface {
	createUser(digestUID string) (*User, error)
	getUserByDigestUID(digestUID string) (*User, error)
	decrementRestCount(user *User) error
}

func (db *DB) createUser(digestUID string) (*User, error) {
	createdAt := time.Now()

	user := User{
		DigestUID: digestUID,
		RestCount: 4,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("db.Create failed: %w", err)
	}

	return &user, nil
}

func (db *DB) getUserByDigestUID(digestUID string) (*User, error) {
	var user User

	if err := db.Where("digest_uid", digestUID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("db.First failed: %w", err)
	}

	return &user, nil
}

func (db *DB) decrementRestCount(user *User) error {
	restCount := user.RestCount

	if restCount == 1 {
		restCount = 4
	} else {
		restCount -= 1
	}

	if err := db.Model(user).Update("rest_count", restCount).Error; err != nil {
		return fmt.Errorf("db.Update failed: %w", err)
	}

	return nil
}
