SET NAMES utf8;
SET time_zone = '+07:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP DATABASE IF EXISTS `myblog`;
CREATE DATABASE `myblog` /*!40100 DEFAULT CHARACTER SET latin1 */;
USE `myblog`;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `password_hash` varchar(255) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `role` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `users` (`id`, `username`, `name`, `password_hash`, `email`, `role`) VALUES
(1,	'hiennm', 'Hien Nguyen', 'password', 'mr.nguyenhien1994@gmail.com', 'admin'),
(2,	'alice', 'Alice', 'password', 'alice@gmail.com', 'user'),
(3,	'bob', 'Bob', 'password', 'bob@gmail.com', 'user');

DROP TABLE IF EXISTS `posts`;
CREATE TABLE `posts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `post_details` varchar(1024) DEFAULT NULL,
  `owner_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  CONSTRAINT `posts_ibfk_1` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

INSERT INTO `posts` (`id`, `title`, `post_details`, `owner_id`) VALUES
(1,	'Lorem Ipsum 1', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam sed massa neque. Ut augue magna, accumsan sed tortor et, aliquam pulvinar leo. Aenean fringilla bibendum dolor vel condimentum. Integer fringilla elit quis lorem hendrerit, eu pretium tellus convallis. Cras malesuada eget lectus id fermentum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam pellentesque consequat sollicitudin. In vel luctus massa. Aenean aliquet volutpat tempus. Nam sed massa nulla.',
2),
(2,	'Lorem Ipsum 2', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam sed massa neque. Ut augue magna, accumsan sed tortor et, aliquam pulvinar leo. Aenean fringilla bibendum dolor vel condimentum. Integer fringilla elit quis lorem hendrerit, eu pretium tellus convallis. Cras malesuada eget lectus id fermentum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam pellentesque consequat sollicitudin. In vel luctus massa. Aenean aliquet volutpat tempus. Nam sed massa nulla.',
3);
