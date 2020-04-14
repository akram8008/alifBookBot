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

	admin,err = InfoUserDB(db,admin)
	if err!=nil {
		log.Fatal("Can not check main admin for exists by error:",err)
	}

	if admin.Role!="admin" {
		err := InsertUser (db,admin)
		if err!=nil {
			log.Fatal("Can not insert main admin")
		}
	}

	return db
}


func InfoUserDB(db *sql.DB, user betypes.User) (betypes.User,error) {
    log.Println("Checking the user: ",user, " in dataBase")
 	row, err := db.Query(userExists, user.ChatId)
	if err!=nil {
		return betypes.User{}, err
	}

	defer row.Close()
 	if row.Next() {
		err = row.Scan(&user.Id, &user.ChatId, &user.FirstName, &user.Phone, &user.Role)
		if err != nil {
			return user, err
		}
		return user, nil
	}
	return  user,nil
}


func InsertUser (db *sql.DB,user betypes.User) error {
	 _, err := db.Exec(insertNewUser,user.ChatId, user.FirstName, user.Phone, user.Role)
	 return err
}

func UpdateUser (db *sql.DB,user betypes.User) error {
	_, err := db.Exec(updateUser,user.ChatId, user.FirstName, user.Phone, user.Role, user.ChatId)
	return err
}
