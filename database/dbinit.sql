DROP DATABASE IF EXISTS `redpendb`;
CREATE DATABASE `redpendb`;

USE `redpendb`;

CREATE TABLE users (
	user_id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    email VARCHAR(255) UNIQUE, 
    password_hash VARCHAR(128),
    reputation INT NOT NULL DEFAULT 0,
    PRIMARY KEY(user_id)
);

USE `redpendb`;

CREATE TABLE session_tokens (
    user_id INT NOT NULL, 
    session_token VARCHAR(64) NOT NULL UNIQUE,
    PRIMARY KEY(user_id)
);

USE `redpendb`;

CREATE TABLE movie_ranks (
    imdb_id VARCHAR(15) NOT NULL,
    user_id INT NOT NULL UNIQUE,
    tier ENUM('S', 'A', 'B', 'C', 'D', 'F'),
    PRIMARY KEY(imdb_id)
);
