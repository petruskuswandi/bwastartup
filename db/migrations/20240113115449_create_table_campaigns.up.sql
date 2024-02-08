CREATE TABLE IF NOT EXISTS `campaigns` (
    `id` CHAR(36) NOT NULL PRIMARY KEY,
    `user_id` CHAR(36) DEFAULT NULL,
    `name` VARCHAR(255) DEFAULT NULL,
    `short_description` VARCHAR(255) DEFAULT NULL,
    `description` TEXT,
    `perks` TEXT,
    `backer_count` INT(11) DEFAULT NULL,
    `goal_amount` INT(11) DEFAULT NULL,
    `current_amount` INT(11) DEFAULT NULL,
    `slug` VARCHAR(255) DEFAULT NULL,
    `created_at` DATETIME DEFAULT NULL,
    `updated_at` DATETIME DEFAULT NULL,
    `deleted_at` DATETIME,
    INDEX `idx_campaign_user_id` (`user_id`),
    CONSTRAINT `fk_campaign_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)
);