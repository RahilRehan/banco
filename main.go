package main

import (
	"banko/db"
	"log"
)

func main(){
	dbh := db.NewDBHandler("postgres://bancoadmin:supersecret@localhost:5432/bancodb?sslmode=disable")
	gdb, err := dbh.GetDB()
	if err!=nil{
		log.Fatalln(err)
	}

	db.CreateAccount(gdb)
	db.GetAccount(gdb)
}
