-- SQL in this section is executed when migration is rolled back.

-- name: remove-queues
DROP TABLE IF EXISTS queues;