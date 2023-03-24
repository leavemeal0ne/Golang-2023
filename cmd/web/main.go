package main

import (
	"github.com/alexedwards/scs/v2"
	config "github.com/leavemeal0ne/Golang-2023/internal/config"
	"github.com/leavemeal0ne/Golang-2023/internal/database"
	"github.com/leavemeal0ne/Golang-2023/internal/driver"
	"github.com/leavemeal0ne/Golang-2023/internal/helpers"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"log"
	"net/http"
	"time"
)

var app config.Config
var session *scs.SessionManager
var template_data *models.TemplateData

func main() {

	app = config.Config{}
	template_data = &models.TemplateData{}
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	app.Session = session

	tc, err := createTemplateCache()
	if err != nil {
		log.Fatal("cannot create temple cache", err)
	}
	app.TemplateCache = tc

	//connection to database
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go user=postgres password=1209")

	if err != nil {
		log.Fatal(err)
	}
	getDatabase(&database.DbRepo{DB: db})
	helpers.GetDatabase(&database.DbRepo{DB: db})
	app.BD = db
	defer db.SQL.Close()

	helpers.HelpersConfig(&app)

	log.Println("Start: http://127.0.0.1:4000")
	err = http.ListenAndServe(":4000", routes(&app))
	if err != nil {
		log.Fatal("Start server Error")
	}

}
