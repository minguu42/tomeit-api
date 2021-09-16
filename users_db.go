package tomeit

import (
	"fmt"
	"time"
)

func (db *DB) createUser(digestUID string) (*user, error) {
	createdAt := time.Now()

	const q = `INSERT INTO users (digest_uid, created_at) VALUES (?, ?)`

	r, err := db.Exec(q, digestUID, createdAt)
	if err != nil {
		return nil, fmt.Errorf("db.Exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("result.lastInsertId failed: %w", err)
	}

	u := user{
		id:        id,
		digestUID: digestUID,
		restCount: 4,
		createdAt: createdAt,
	}
	return &u, nil
}

func (db *DB) getUserByDigestUID(digestUID string) (*user, error) {
	const q = `SELECT id, digest_uid, rest_count, created_at FROM users WHERE digest_uid = ?`

	var u user
	if err := db.QueryRow(q, digestUID).Scan(&u.id, &u.digestUID, &u.restCount, &u.createdAt); err != nil {
		return nil, fmt.Errorf("db.QueryRow failed: %w", err)
	}

	return &u, nil
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
