-- SQL in this section is executed when migration is rolled back.

-- name: remove-user-accounts
DROP TABLE IF EXISTS user_accounts;