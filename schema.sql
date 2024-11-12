CREATE TABLE video (
    id   UUID PRIMARY KEY,
    playlist VARCHAR(255) NOT NULL,
    videoID VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    thumbnail VARCHAR(255) NOT NULL,
    duration VARCHAR(255) NOT NULL
);