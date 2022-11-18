
CREATE TABLE IF NOT EXISTS service_table (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(250) NOT NULL,
    price INT NOT NULL,
    availability boolean default TRUE
);

CREATE TABLE IF NOT EXISTS balance_table (
    id SERIAL PRIMARY KEY,
    user_id INT DEFAULT NULL,
    total INT,
    reserve INT default NULL,
    order_number VARCHAR(250) NOT NULL,
    service_id INT default null,
    created_at timestamp default NULL,
    status VARCHAR(100) default NULL
);

CREATE TABLE IF NOT EXISTS user_table (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(250) NOT NULL,
    last_name VARCHAR(250) NOT NULL,
    balance_id INT DEFAULT NULL
);

ALTER TABLE user_table
    ADD CONSTRAINT fk_links_user_table_balance_table
    FOREIGN KEY (balance_id)
    REFERENCES balance_table (id);

ALTER TABLE balance_table
    ADD CONSTRAINT fk_links_balance_table_service_table
    FOREIGN KEY (service_id)
    REFERENCES service_table(id);

ALTER TABLE balance_table
    ADD CONSTRAINT fk_links_balance_table_user_table
    FOREIGN KEY (user_id)
    REFERENCES user_table (id);
