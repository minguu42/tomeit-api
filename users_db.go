package tomeit

import "fmt"

func (db *DB) createUser(digestUID string) (*user, error) {
	const q = `INSERT INTO users (digest_uid) VALUES (?)`

	r, err := (*db).Exec(q, digestUID)
	if err != nil {
		return nil, fmt.Errorf("exec failed: %w", err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("lastInsertId failed: %w", err)
	}

	u := user{
		id:        id,
		digestUID: digestUID,
		restCount: 4,
	}
	return &u, nil
}

func (db *DB) getUserByDigestUID(digestUID string) (*user, error) {
	const q = `SELECT * FROM users WHERE digest_uid = ?`

	var u user
	if err := db.QueryRow(q, digestUID).Scan(&u.id, &u.digestUID, &u.restCount); err != nil {
		return nil, fmt.Errorf("queryRow failed: %w", err)
	}

	return &u, nil
}
