
USE `teirxdb`;

CREATE TABLE session_tokens (
    user_id INT NOT NULL, 
    session_token VARCHAR(64) NOT NULL UNIQUE,
    PRIMARY KEY(user_id)
);
