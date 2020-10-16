package main

import (
	"assignment/conf"
	_ "assignment/docs"
	"assignment/logger"
	"assignment/pg"
	"assignment/repo"
	"assignment/server"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/swaggo/http-swagger"
	"net/http"
)

// @title Swagger Example API
// @version 1.0
// @description This is an assignment
func main() {
	config := conf.GetConfig()
	spew.Dump(config)
	logger.InitLogger(config.Log)
	db := pg.InitDB(config.DB)
	repo.Init(db)
	router := server.NewRouter(config.Server)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition"
	))
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), router)
}
