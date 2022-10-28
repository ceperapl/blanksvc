CREATE TABLE tasks
(
    id              uuid            PRIMARY KEY,
    name            varchar(255)    NOT NULL,
    description     TEXT            DEFAULT NULL,
    deadline        date            NOT NULL,
    completed_at    timestamptz      DEFAULT NULL, 
    created_at      timestamptz      DEFAULT current_timestamp,
    updated_at      timestamptz      DEFAULT NULL
);
