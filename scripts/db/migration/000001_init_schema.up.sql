CREATE TABLE `accounts`
(
    `id`          BIGINT AUTO_INCREMENT PRIMARY KEY,
    `email`       varchar(255) NOT NULL UNIQUE,
    `password`    varchar(255) NOT NULL,
    `name`        varchar(255) NOT NULL,
    `is_valid`    boolean      NOT NULL DEFAULT false,
    `verified_at` timestamp    NULL DEFAULT NULL,
    `created_at`  timestamp    NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    `updated_at`  timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE `account_valid_temp`
(
    `id`         BIGINT AUTO_INCREMENT PRIMARY KEY,
    `account_id` bigint       NOT NULL,
    `code`       varchar(255) NOT NULL UNIQUE,
    `created_at` timestamp    NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    `expired_at` timestamp    NOT NULL
);

CREATE TABLE `products`
(
    `id`         BIGINT AUTO_INCREMENT PRIMARY KEY,
    `name`       varchar(255) NOT NULL,
    `price`      decimal      NOT NULL,
    `created_at` timestamp    NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    `updated_at` timestamp    NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE `recommend_products`
(
    `id`         BIGINT AUTO_INCREMENT PRIMARY KEY,
    `product_id` bigint    NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP),
    `updated_at` timestamp NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE INDEX `accounts_index_0` ON `accounts` (`email`);

ALTER TABLE `account_valid_temp`
    ADD FOREIGN KEY (`account_id`) REFERENCES `accounts` (`id`);

ALTER TABLE `recommend_products`
    ADD FOREIGN KEY (`product_id`) REFERENCES `products` (`id`);

# init
INSERT INTO products (name, price, created_at, updated_at)
VALUES ('test_product_1', 10.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO products (name, price, created_at, updated_at)
VALUES ('test_product_2', 20.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO products (name, price, created_at, updated_at)
VALUES ('test_product_3', 30.00, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO recommend_products (product_id)
VALUES (1);
INSERT INTO recommend_products (product_id)
VALUES (2);