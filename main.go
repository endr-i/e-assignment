package main

import (
	"assignment/conf"
	"assignment/pg"
	"assignment/repo"
	"assignment/server"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

func main() {
	config := conf.GetConfig()
	spew.Dump(config)
	db := pg.InitDB(config.DB)
	repo.Init(db)
	router := server.NewRouter(config.Server)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Server.Port), router)
}
