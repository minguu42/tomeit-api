package tomeit

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
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
