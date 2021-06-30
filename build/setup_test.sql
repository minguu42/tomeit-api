SET CHARACTER SET UTF8;

CREATE TABLE IF NOT EXISTS users (
    id              INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    digest_uid      CHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id         INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id    INT          NOT NULL,
    name       VARCHAR(120) NOT NULL,
    priority   INT          DEFAULT 0 NOT NULL,
    deadline   DATE         DEFAULT ('0001-01-01') NOT NULL,
    is_done    BOOLEAN      DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS pomodoro_logs (
    id         INT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id    INT       NOT NULL,
    task_id    INT       NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

INSERT INTO users (digest_uid) VALUES ('a2c4ba85c41f186283948b1a54efacea04cb2d3f54a88d5826a7e6a917b28c5a');
INSERT INTO tasks (user_id, name, priority, deadline, is_done) VALUES (1, 'タスク1', 0, '2021-06-30', false),
                                                                      (1, 'タスク2', 1, '2021-07-01', false),
                                                                      (1, '完了済みタスク1', 0, '2021-06-01', true),
                                                                      (1, '完了済みタスク2', 2, '2021-06-10', true);
