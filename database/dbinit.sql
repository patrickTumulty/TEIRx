DROP DATABASE IF EXISTS `reviewdb`;
CREATE DATABASE `reviewdb`;
USE `reviewdb`;

CREATE TABLE users (
	user_id INT NOT NULL AUTO_INCREMENT,
    firstname VARCHAR(50),
    lastname VARCHAR(50),
    date_of_birth DATE,
    zipcode VARCHAR(5),
    email VARCHAR(254)
    PRIMARY KEY(user_id)
);
