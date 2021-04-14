package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/MinhL4m/blogs/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(username, password, host, port, dbname, dbtype string) {
	var err error
	var DBURL string

	dbtype = strings.ToLower(dbtype)

	switch dbtype {
	case "mysql":
		DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)

	default:
		DBURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, username, dbname, password)

	}

	a.DB, err = gorm.Open(dbtype, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", dbtype)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", dbtype)
	}

	// migration
	a.DB.Debug().AutoMigrate(&models.User{}, &models.Post{})
	a.Router = mux.NewRouter()

	// TODO:
	// a.initializeRoute()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
