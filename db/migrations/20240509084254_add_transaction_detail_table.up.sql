CREATE TABLE IF NOT EXISTS transaction_detail(
    transaction_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL,

    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
    ON UPDATE CASCADE ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
    ON UPDATE CASCADE ON DELETE CASCADE
);