CREATE TABLE IF NOT EXISTS users
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    chatid    INTEGER,
    firstname TEXT NOT NULL,
    phone     TEXT NOT NULL,
    role      TEXT NOT NULL
);


insertNewUser = `INSERT INTO users(chatid, firstname, phone, role)  VALUES (?, ?, ?, ?);`

UPDATE users SET chatid=?, firstname=?, phone=?, role=?  WHERE chatid = ?;


SELECT * FROM users WHERE chatid=5;



UPDATE firstname, phone, role FROM users WHERE chatid=;

drop table users;




UPDATE users SET chatid=?, firstname=?, phone=?, role=?  WHERE chatid = ?;

