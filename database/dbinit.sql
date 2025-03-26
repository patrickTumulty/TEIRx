DROP DATABASE IF EXISTS `teirxdb`;
CREATE DATABASE `teirxdb`;
USE `teirxdb`;

CREATE TABLE users (
	user_id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    email VARCHAR(254) UNIQUE, 
    password_hash VARCHAR(128),
    reputation INT NOT NULL,
    PRIMARY KEY(user_id)
);
