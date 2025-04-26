CREATE DATABASE IF NOT EXISTS gophersignal;
USE gophersignal;

CREATE TABLE IF NOT EXISTS articles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    hn_id INT NOT NULL DEFAULT 0, 
    title VARCHAR(255) NOT NULL,
    link VARCHAR(512) NOT NULL,
    article_rank INT NOT NULL, 
    content TEXT,
    summary VARCHAR(2000),
    source VARCHAR(100) NOT NULL,
    upvotes INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    comment_link VARCHAR(255),
    flagged BOOLEAN NOT NULL DEFAULT FALSE,
    dead BOOLEAN NOT NULL DEFAULT FALSE,
    dupe BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
