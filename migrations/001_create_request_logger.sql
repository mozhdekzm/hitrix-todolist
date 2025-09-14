-- +migrate Up
CREATE TABLE request_logger (
                                ID BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
                                URL VARCHAR(255),
                                UserID BIGINT UNSIGNED,
                                AppName VARCHAR(255) NOT NULL,
                                Request MEDIUMBLOB,
                                RequestText TEXT,
                                Response MEDIUMBLOB,
                                ResponseText TEXT,
                                Log MEDIUMBLOB,
                                Status INT,
                                CreatedAt DATETIME NOT NULL,
                                INDEX idx_created_at (CreatedAt)
);
-- +migrate Down
DROP TABLE IF EXISTS request_logger;
