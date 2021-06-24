package tomeit

func getUserByDigestUID(digestUID string) (User, error) {
	const q = `SELECT * FROM users WHERE digest_uid = ?`

	var u User
	if err := db.QueryRow(q, digestUID).Scan(&u.id, &u.digestUID); err != nil {
		return User{}, err
	}

	return u, nil
}

func createUser(digestUID string) (User, error) {
	const q = `
INSERT INTO users (digest_uid)
VALUES (?);
`
	r, err := db.Exec(q, digestUID)
	if err != nil {
		return User{}, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return User{}, err
	}

	u := User{
		id:        id,
		digestUID: digestUID,
	}

	return u, nil
}
