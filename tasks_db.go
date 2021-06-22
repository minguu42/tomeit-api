package tomeit

func insertTask(userId int, name string, priority int, deadline string) (Task, error) {
	const query = `INSERT INTO tasks (user_id, name, priority, deadline) VALUES ($1, $2, $3, $4) RETURNING id;`
	// TODO: 作成途中, DB を触れるようにする
	return Task{}, nil
}
