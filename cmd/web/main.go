package main

import (
	"database/sql"
	"fmt"
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
var EnvConfig EnvConfigModel

type EnvConfigModel struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	SSLMode        string `mapstructure:"SSL_MODE"`
}

func LoadConfig(filePath string) (err error) {
	viper.SetConfigType("env")
	viper.SetConfigFile(filePath)

	viper.AutomaticEnv()

	if viper.ReadInConfig() != nil {
		return
	}

	return viper.Unmarshal(&EnvConfig)
}

func main() {

	err := LoadConfig(".env")
	if err != nil {
		log.Fatalln("Failed to load environment variables!", err.Error())
	}

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
	db, err := driver.ConnectSQL(fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", EnvConfig.DBHost,
		EnvConfig.DBPort, EnvConfig.DBName, EnvConfig.DBUserName, EnvConfig.DBUserPassword))

	log.Println(fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", EnvConfig.DBHost,
		EnvConfig.DBPort, EnvConfig.DBName, EnvConfig.DBUserName, EnvConfig.DBUserPassword))

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

	err = http.ListenAndServe(":8000", routes(&app))
	if err != nil {
		log.Fatal("Start server Error")
	}

}
