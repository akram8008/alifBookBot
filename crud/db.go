package dataBase

import (
	"alifLibrary/betypes"
	"database/sql"
	"log"
)



func Connect () *sql.DB {
	log.Print("connecting to database")

	db, err := sql.Open("sqlite3", "sqLite.DB")
	if err != nil {
		log.Fatalf("can't open sqlite3: %v", err)
	}

	log.Print("recreate tables if not exist in database")
	_, err = db.Exec(createTable)
	if err!=nil {
		log.Fatal("Can not recreate main tables")
	}

	log.Print("add main admin if he is not registered in database")
	admin := betypes.User{
		ChatId: betypes.AdminChatId,
		FirstName: betypes.AdminFirstName,
		Phone:betypes.AdminPhone,
		Role:betypes.AdminRole,
	}

	ok,err := IsUserExist(db,admin)
	if err!=nil {
		log.Fatal("Can not check main admin for exists")
	}

	if !ok {
		err := InsertUser (db,admin)
		if err!=nil {
			log.Fatal("Can not insert main admin")
		}
	}

	return db
}


func IsUserExist (db *sql.DB, user betypes.User) (bool,error) {
	var id int
	err := db.QueryRow(userExists, user.Phone).Scan(&id)
	if err == nil {
		return true,nil
	}
	if err == sql.ErrNoRows {
		return false,nil
	}
	return false,err
}


func InsertUser (db *sql.DB,user betypes.User) error {
	 _, err := db.Exec(insertNewUser,user.FirstName,user.Phone,user.Role)
	 return err
}