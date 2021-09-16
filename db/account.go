package db

import (
	"fmt"
	"gorm.io/gorm"
)

type Account struct{
	Id int
	Owner string
	Balance string
	Currency string
}

func CreateAccount(db *gorm.DB){
	ac := Account{
		Owner: "Rehan",
		Balance: "4321",
		Currency: "dollar",
	}

	create := db.Create(&ac)
	fmt.Println(create.RowsAffected)
}

func GetAccount(db *gorm.DB){
	var ac []Account
	db.Find(&ac)
	fmt.Printf("%+v", ac)
}
