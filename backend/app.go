package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(username, password, dbname, dbtype string) {

}

func (a *App) Run(addr string) {
	http.ListenAndServe(addr, a.Router)
}
