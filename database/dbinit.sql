DROP DATABASE IF EXISTS `teirxdb`;
CREATE DATABASE `teirxdb`;
USE `teirxdb`;

CREATE TABLE users (
	user_id INT NOT NULL AUTO_INCREMENT,
    firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    date_of_birth DATE NOT NULL,
    zipcode VARCHAR(5) NOT NULL,
    email VARCHAR(254) UNIQUE,
    PRIMARY KEY(user_id)
);
