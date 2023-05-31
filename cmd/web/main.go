package main

import (
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"github.com/leavemeal0ne/Golang-2023/internal/config"
	"github.com/leavemeal0ne/Golang-2023/internal/database"
	"github.com/leavemeal0ne/Golang-2023/internal/driver"
	"github.com/leavemeal0ne/Golang-2023/internal/helpers"
	"github.com/leavemeal0ne/Golang-2023/internal/models"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

var app config.Config
var session *scs.SessionManager
var template_data *models.TemplateData

func main() {

	vi := viper.New()
	vi.SetConfigName("config.json")
	vi.SetConfigType("json")

	//host := vi.GetString("host")
	//port := vi.GetInt("port")
	//dbName := vi.GetString("dbname")
	//user := vi.GetString("user")
	//password := vi.GetInt("password")

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
		log.Println("cannot create temple cache", err)
	}
	app.TemplateCache = tc

	//connection to database
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=go user=postgres password=1209")

	getDatabase(&database.DbRepo{DB: db})
	helpers.GetDatabase(&database.DbRepo{DB: db})
	app.BD = db
	defer func(SQL *sql.DB) {
		err := SQL.Close()
		if err != nil {

		}
	}(db.SQL)

	helpers.HelpersConfig(&app)

	if err != nil {
		log.Println(err)
		log.Fatal("Start server error")
	}

	log.Println("Start: http://127.0.0.1:4000")
	err = http.ListenAndServe(":4000", routes(&app))
	if err != nil {
		log.Fatal("Start server Error")
	}

}
