-- SQL in this section is executed when migration is rolled back.

-- name: remove-customers
DROP TABLE IF EXISTS customers;