CREATE DATABASE IF NOT EXISTS mydb;

USE mydb;

CREATE TABLE IF NOT EXISTS compromised_emails (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    compromised BOOLEAN NOT NULL DEFAULT FALSE,
    compromised_at TIMESTAMP NULL DEFAULT NULL
);

INSERT INTO compromised_emails (email, compromised, compromised_at) VALUES
('test1@example.com', TRUE, NOW()),
('test2@example.com', FALSE, NULL),
('test3@example.com', TRUE, NOW() - INTERVAL 7 DAY),
('user4@example.com', TRUE, NOW() - INTERVAL 30 DAY);