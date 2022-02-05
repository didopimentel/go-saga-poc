CREATE TABLE orders (
    id bigserial PRIMARY KEY,
    amount bigint NOT NULL
);

CREATE TABLE payments (
    id bigserial PRIMARY KEY,
    order_id bigint NOT NULL
);

CREATE TABLE deliveries (
  id bigserial PRIMARY KEY,
  order_id bigint NOT NULL
)