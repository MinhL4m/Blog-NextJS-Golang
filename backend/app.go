package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB    *sql.DB
}

func (a *App) Initialize(username, password, dbname, dbtype string) {

}

func (a *App) Run(addr string) {
	http.ListenAndServe(addr, a.Router)
}
