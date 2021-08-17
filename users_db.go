package tomeit

import (
	"fmt"
	"time"
)

func (db *DB) createUser(digestUID string) (*user, error) {
	createdAt := time.Now()

	const q = `INSERT INTO users (digest_uid, created_at) VALUES (?, ?)`

	r, err := db.Exec(q, digestUID)
	if err != nil {
		return nil, fmt.Errorf("db.Exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("result.lastInsertId failed: %w", err)
	}

	u := user{
		id:            id,
		digestUID:     digestUID,
		nextRestCount: 4,
		createdAt:     createdAt,
	}
	return &u, nil
}

func (db *DB) getUserByDigestUID(digestUID string) (*user, error) {
	const q = `SELECT * FROM users WHERE digest_uid = ?`

	var u user
	if err := db.QueryRow(q, digestUID).Scan(&u.id, &u.digestUID, &u.nextRestCount, &u.createdAt); err != nil {
		return nil, fmt.Errorf("db.QueryRow failed: %w", err)
	}

	return &u, nil
}

func (db *DB) decrementNextRestCount(user *user) error {
	nextRestCount := user.nextRestCount

	if nextRestCount == 1 {
		nextRestCount = 4
	} else {
		nextRestCount -= 1
	}

	const q = `UPDATE users SET next_rest_count = ? WHERE id = ?`

	if _, err := db.Exec(q, nextRestCount, user.id); err != nil {
		return fmt.Errorf("db.Exec failed: %w", err)
	}

	return nil
}
