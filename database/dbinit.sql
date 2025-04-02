DROP DATABASE IF EXISTS `teirxdb`;
CREATE DATABASE `teirxdb`;

USE `teirxdb`;

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

USE `teirxdb`;

CREATE TABLE session_tokens (
    user_id INT NOT NULL, 
    session_token VARCHAR(64) NOT NULL UNIQUE,
    PRIMARY KEY(user_id)
);

USE `teirxdb`;

CREATE TABLE movie_ranks (
    imdb_id VARCHAR(15) NOT NULL,
    s_teir INT NOT NULL DEFAULT 0,
    a_teir INT NOT NULL DEFAULT 0,
    b_teir INT NOT NULL DEFAULT 0,
    c_teir INT NOT NULL DEFAULT 0,
    d_teir INT NOT NULL DEFAULT 0,
    f_teir INT NOT NULL DEFAULT 0,
    PRIMARY KEY(imdb_id)
);
