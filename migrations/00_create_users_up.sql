-- SQL in this section is executed when migration is applied.

-- name: create-user-accounts
CREATE TABLE IF NOT EXISTS user_accounts
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    username        VARCHAR(255)   NOT NULL,
    password        VARCHAR(255)   NOT NULL,
    is_admin        BOOLEAN        DEFAULT FALSE,
    created_at      TIMESTAMP      DEFAULT CURRENT_TIMESTAMP(),
    last_login_at   TIMESTAMP      NULL,
    updated_at      TIMESTAMP      NULL,      
    PRIMARY KEY(id)
);

-- name: create-username-index
CREATE UNIQUE INDEX user_accounts_username_index ON user_accounts(username);