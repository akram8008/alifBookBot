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

	admin,err = IsUserExist(db,admin)
	if err!=nil {
		log.Fatal("Can not check main admin for exists")
	}

	if admin.Role!="admin" {
		err := InsertUser (db,admin)
		if err!=nil {
			log.Fatal("Can not insert main admin")
		}
	}

	return db
}


func IsUserExist (db *sql.DB, checkUser betypes.User) (betypes.User,error) {
	user := betypes.User{}
	err := db.QueryRow(userExists, checkUser.ChatId).Scan(&user.Id,&user.ChatId,&user.FirstName,&user.Phone,&user.Role)

	if err == nil {
		return user,nil
	}
	if err == sql.ErrNoRows {
		return betypes.User{},nil
	}
	return betypes.User{},err
}


func InsertUser (db *sql.DB,user betypes.User) error {
	 _, err := db.Exec(insertNewUser,user.ChatId, user.FirstName, user.Phone, user.Role)
	 return err
}

func UpdateUser (db *sql.DB,user betypes.User) error {
	_, err := db.Exec(updateUser,user.ChatId, user.FirstName, user.Phone, user.Role, user.ChatId)
	return err
}
