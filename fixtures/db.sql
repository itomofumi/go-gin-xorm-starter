
DROP DATABASE IF EXISTS `go-gin-xorm-starter`;

CREATE DATABASE `go-gin-xorm-starter`;

USE `go-gin-xorm-starter`;

/*
  CREATE TABLES
*/

CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `email` varchar(120) NOT NULL,
  `email_verified` tinyint(1) NOT NULL,
  `display_name` varchar(50) DEFAULT NULL,
  `about` text,
  `avatar_url` text,
  `last_login_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_users_pk` (`id`),
  KEY `IDX_users_mail` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `fruits` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `is_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `name` varchar(255) NOT NULL,
  `price` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `IDX_fruits_pk` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


/*
  INSERT DATA
*/

INSERT INTO `users` (`is_deleted`, `is_enabled`, `created_at`, `updated_at`, `email`, `email_verified`, `display_name`, `about`, `avatar_url`, `last_login_at`)
VALUES
  (0,1,'2018-01-01 00:00:00','2018-01-01 00:00:00','test@example.com',1,'テストユーザー','テストユーザーです','https://s3-ap-northeast-1.amazonaws.com/gemcook.com/assets/images/gemo_houseki.png','2018-01-01 00:00:00');

INSERT INTO `fruits` (`is_deleted`, `is_enabled`, `created_at`, `updated_at`, `name`, `price`)
VALUES 
  (0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Apple', 112),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Pear', 245),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Banana', 60),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Orange', 80),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Kiwi', 106),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Strawberry', 350),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Grape', 400),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Grapefruit', 150),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Pineapple', 200),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Cherry', 140),
	(0, 1, '2018-01-01 00:00:00', '2018-01-01 00:00:00', 'Mango', 199);
