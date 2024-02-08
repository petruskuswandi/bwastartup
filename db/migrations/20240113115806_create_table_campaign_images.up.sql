CREATE TABLE IF NOT EXISTS `campaign_images` (
    `id` CHAR(36) NOT NULL PRIMARY KEY,
    `campaign_id` CHAR(36) DEFAULT NULL,
    `file_name` VARCHAR(255) DEFAULT NULL,
    `is_primary` TINYINT(1) DEFAULT NULL,
    `created_at` DATETIME DEFAULT NULL,
    `updated_at` DATETIME DEFAULT NULL,
    `deleted_at` DATETIME,
    INDEX `idx_campaign_image_campaign_id` (`campaign_id`),
    CONSTRAINT `fk_campaign_image_campaign_id` FOREIGN KEY (`campaign_id`) REFERENCES `campaigns`(`id`)
);