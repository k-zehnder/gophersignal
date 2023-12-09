USE gopher_api;

CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    text VARCHAR(255) NOT NULL,
    tags TEXT,
    due DATETIME
);

INSERT INTO tasks (text, tags, due) VALUES ('Sample Task', '["tag1","tag2"]', NOW());
