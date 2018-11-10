-- SQL in this section is executed when migration is applied.

-- name: create-customers
CREATE TABLE IF NOT EXISTS customers
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    msisdn          VARCHAR(255)   NOT NULL,
    ticket          VARCHAR(255)   NOT NULL,
    queue_id        INT            NOT NULL,
    created_at      DATETIME       DEFAULT NOW(),
    served_at       DATETIME       NULL,     
    PRIMARY KEY(id),
    CONSTRAINT fk_customers_queue_id  FOREIGN KEY  (queue_id)     REFERENCES queues(id)
);

-- name: create-customers-msisdn_index
CREATE INDEX customers_msisdn_index ON customers(msisdn);