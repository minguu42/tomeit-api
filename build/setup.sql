CREATE TABLE IF NOT EXISTS users (
    id              INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    digest_id_token VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id         INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id    INT          NOT NULL REFERENCES users(id),
    name       VARCHAR(120) NOT NULL,
    priority   INT          NOT NULL DEFAULT 0,
    deadline   DATE         NOT NULL DEFAULT ('0001-01-01'),
    is_done    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pomodoroLogs (
    id         INT      NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id    INT      NOT NULL REFERENCES users(id),
    task_id    INT      NOT NULL REFERENCES tasks(id),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);