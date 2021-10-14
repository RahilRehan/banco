package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/RahilRehan/banco/api"
	db "github.com/RahilRehan/banco/db/sqlc"
	"github.com/RahilRehan/banco/db/util"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("Cannot read config: ", err)
	}
	dbSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB_USER,
		os.Getenv("DB_PASSWORD"),
		cfg.DB_HOST,
		cfg.DB_PORT,
		cfg.DB_NAME,
		cfg.SSL_MODE,
	)

	conn, err := sql.Open(cfg.DRIVER_NAME, dbSource)
	if err != nil {
		log.Fatalln("Cannot connect to DB ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(*cfg, store)
	if err != nil {
		log.Fatalln("Cannot start server ", err)
	}
	err = server.Start(cfg.SERVER_ADDRESS)
	if err != nil {
		log.Fatalln("Cannot start server ", err)
	}
}
