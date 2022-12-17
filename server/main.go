package main

import (
	"database/sql"
	"log"

	"github.com/brkss/vanillefraise2/api"
	db "github.com/brkss/vanillefraise2/db/sqlc"
	"github.com/brkss/vanillefraise2/utils"
	_ "github.com/lib/pq"
)

func main() {

	config, err := utils.LoadConfig()

	if err != nil {
		log.Fatal("cannot load config :", err)
	}

	con, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to databse :", err)
	}

	source := db.New(con)
	server := api.NewServer(source)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("connot connect to server :", err)
	}
}
