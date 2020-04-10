CREATE TABLE IF NOT EXISTS users(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    firstname TEXT NOT NULL,
    phone     TEXT NOT NULL,
    role      TEXT NOT NULL
);