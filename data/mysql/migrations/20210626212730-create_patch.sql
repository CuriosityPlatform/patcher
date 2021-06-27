-- +migrate Up
CREATE TABLE patch
(
    `patch_id` binary(16) NOT NULL,
    `project` varchar(255) NOT NULL,
    `applied` tinyint(1) NOT NULL,
    `content` longtext NOT NULL,
    `author` varchar(255) NOT NULL,
    `device` varchar(255) NOT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`patch_id`),
    INDEX `patch_id_index` (`patch_id`)
);

-- +migrate Down
DROP TABLE patch;