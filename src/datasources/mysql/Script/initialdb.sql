CREATE DATABASE IF NOT EXISTS lms_db;

USE lms_db;

CREATE TABLE `user` (
                        `id` INT NOT NULL AUTO_INCREMENT COMMENT 'Primary key.',
                        `domain` varchar(255) NOT NULL COMMENT 'Domain of the company the user belongs to',
                        `user_id` varchar(255) NOT NULL COMMENT 'Id of a user. Will be unique per domain.',
                        `metadata` json DEFAULT NULL COMMENT 'Metadata JSON. Will store anything about a user.',
                        `status` varchar(255) DEFAULT NULL COMMENT 'Special kind of meta.',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY uk_domain_user_id (`domain`,`user_id`)
) ENGINE=InnoDB CHARSET=utf8;

CREATE USER 'md_user'@'%' IDENTIFIED BY '12345678';
GRANT ALL PRIVILEGES ON *.* TO 'md_user'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
