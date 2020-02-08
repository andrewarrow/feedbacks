CREATE TABLE inbox (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    sent_to varchar(255),
    sent_from varchar(255),
    subject varchar(255),
    host varchar(255),
    is_spam int,
    spam_score float,
    body text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    KEY host (host)
) ENGINE InnoDB;
