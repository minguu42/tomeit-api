SET CHARACTER SET UTF8;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS pomodoro_records;

CREATE TABLE IF NOT EXISTS users (
    id              INT       PRIMARY KEY AUTO_INCREMENT,
    digest_uid      CHAR(64)  NOT NULL UNIQUE,
    next_rest_count INT       DEFAULT 4 NOT NULL CHECK ( 1 <= rest_count AND rest_count <= 4 ),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id                       INT          PRIMARY KEY AUTO_INCREMENT,
    user_id                  INT          NOT NULL,
    title                    VARCHAR(120) NOT NULL,
    expected_pomodoro_number INT          DEFAULT 0 NOT NULL CHECK (0 <= priority AND priority <= 3),
    due_on                   DATE         DEFAULT ('0001-01-01') NOT NULL,
    is_completed             BOOLEAN      DEFAULT FALSE NOT NULL,
    completed_at             TIMESTAMP    DEFAULT ('0001-01-01 00:00:00') NOT NULL,
    created_at               TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at               TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS pomodoro_records (
    id           INT       PRIMARY KEY AUTO_INCREMENT,
    user_id      INT       NOT NULL,
    task_id      INT       NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);