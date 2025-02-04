CREATE DATABASE IF NOT EXISTS gophersignal;
USE gophersignal;

CREATE TABLE IF NOT EXISTS articles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    link VARCHAR(512) NOT NULL,
    content TEXT,
    summary VARCHAR(2000),
    source VARCHAR(100) NOT NULL,
    upvotes INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    comment_link TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
