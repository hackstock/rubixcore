-- SQL in this section is executed when migration is applied.

-- name: create-customers
CREATE TABLE IF NOT EXISTS customers
(
    id              INT            NOT NULL     AUTO_INCREMENT,
    msisdn          VARCHAR(255)   NOT NULL,
    ticket          VARCHAR(255)   NOT NULL,
    queue_id        INT            NOT NULL,
    created_at      TIMESTAMP      DEFAULT CURRENT_TIMESTAMP(),
    served_by       INT            NULL, 
    served_at       TIMESTAMP      NULL,     
    PRIMARY KEY(id),
    CONSTRAINT fk_customers_queue_id  FOREIGN KEY  (queue_id)     REFERENCES queues(id),
    CONSTRAINT fk_customers_served_by FOREIGN KEY  (served_by)    REFERENCES user_accounts(id)
);

-- name: create-customers-msisdn_index
CREATE INDEX customers_msisdn_index ON customers(msisdn);