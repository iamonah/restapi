CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY,
    slug TEXT NOT NULL,
    author TEXT NOT NULL,
    body TEXT NOT NULL
);