CREATE TABLE IF NOT EXISTS `transactions` (
    `id` CHAR(36) NOT NULL PRIMARY KEY,
    `campaign_id` CHAR(36) DEFAULT NULL,
    `user_id` CHAR(36) DEFAULT NULL,
    `amount` INT(11) DEFAULT NULL,
    `status` VARCHAR(255) DEFAULT NULL,
    `code` VARCHAR(255) DEFAULT NULL,
    `payment_url` VARCHAR(255),
    `created_at` DATETIME DEFAULT NULL,
    `updated_at` DATETIME DEFAULT NULL,
    `deleted_at` DATETIME,
    INDEX `idx_transaction_campaign_id` (`campaign_id`),
    INDEX `idx_transaction_user_id` (`user_id`),
    CONSTRAINT `fk_transaction_campaign_id` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns`(`id`),
    CONSTRAINT `fk_transaction_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
);