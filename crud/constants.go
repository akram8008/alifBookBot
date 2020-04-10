package dataBase

const (
createTable = `CREATE TABLE IF NOT EXISTS   users (
   id                 INTEGER PRIMARY KEY AUTOINCREMENT,
   chatid             INTEGER,
   firstname          TEXT NOT NULL,	
   phone              TEXT NOT NULL,
   role               TEXT NOT NULL
);
`

insertNewUser = `INSERT INTO users(firstname, phone, role)  VALUES (?, ?, ?);`
userExists = `SELECT id FROM users WHERE phone=?;`
)
