package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RahilRehan/banco/api"
	"github.com/RahilRehan/banco/db"
	"github.com/RahilRehan/banco/db/util"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {

	cfg, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalln("Cannot read config: ", err)
	}
	dbSource := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB_USER,
		viper.GetString("DB_PASSWORD"),
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
	server := api.NewServer(store)

	err = server.Start(cfg.SERVER_ADDRESS)
	if err != nil {
		log.Fatalln("Cannot start server ", err)
	}
}
