-- SQL in this section is executed when migration is applied.

-- name: create-queues
CREATE TABLE IF NOT EXISTS queues
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    name            VARCHAR(255)   NOT NULL,
    description     VARCHAR(255)   NOT NULL,
    is_active       BOOLEAN        DEFAULT FALSE,
    created_at      DATETIME       DEFAULT NOW(),
    updated_at      TIMESTAMP      NULL,      
    PRIMARY KEY(id)
);

-- name: create-queues-name-index
CREATE UNIQUE INDEX queues_name_index ON queues(name);