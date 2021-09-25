DROP TABLE IF EXISTS pomodoros;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users (
    id         INT       PRIMARY KEY AUTO_INCREMENT,
    digest_uid CHAR(64)  NOT NULL UNIQUE,
    rest_count INT       DEFAULT 4 NOT NULL CHECK ( 1 <= rest_count AND rest_count <= 4 ),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS tasks (
    id                       INT          PRIMARY KEY AUTO_INCREMENT,
    user_id                  INT          NOT NULL,
    title                    VARCHAR(120) NOT NULL,
    expected_pomodoro_num    INT          DEFAULT 0 NOT NULL CHECK (0 <= expected_pomodoro_number AND expected_pomodoro_number <= 6),
    due_at                   TIMESTAMP,
    is_completed             BOOLEAN      DEFAULT FALSE NOT NULL,
    completed_at             TIMESTAMP,
    created_at               TIMESTAMP    DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at               TIMESTAMP    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS pomodoros (
    id           INT       PRIMARY KEY AUTO_INCREMENT,
    user_id      INT       NOT NULL,
    task_id      INT       NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);