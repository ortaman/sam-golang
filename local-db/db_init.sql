
CREATE TABLE IF NOT EXISTS customer (
    customer_id SERIAL  PRIMARY KEY,
    name varchar(100) NULL,
    email varchar(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    trans_date Date,
    amount FLOAT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    customer_id INT REFERENCES customer(customer_id)
);
