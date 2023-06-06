package main

import (
	"database/sql"
	"fmt"
	"github.com/leavemeal0ne/Golang-2023/internal/config"
	"github.com/leavemeal0ne/Golang-2023/internal/driver"
	"github.com/leavemeal0ne/Golang-2023/internal/handlers"
	"github.com/leavemeal0ne/Golang-2023/internal/helpers"
	"github.com/leavemeal0ne/Golang-2023/internal/middleware"
	"github.com/leavemeal0ne/Golang-2023/internal/repositories"
	"github.com/leavemeal0ne/Golang-2023/internal/router"
	"github.com/leavemeal0ne/Golang-2023/internal/services"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var app config.Config
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

	app = *config.NewAppConfig()

	middleware.SetSession(app.Session)

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

	defer func(SQL *sql.DB) {
		err := SQL.Close()
		if err != nil {

		}
	}(db.SQL)

	if err != nil {
		log.Println(err)
		log.Fatal("Start server error")
	}

	helper := helpers.NewHelper(&app)

	noteRepository := repositories.NewNoteRepository(db)
	userRepository := repositories.NewUserRepository(db)

	noteService := services.NewNoteService(noteRepository)
	userService := services.NewUserService(userRepository)

	noteHttp := handlers.NewNoteHttpHandler(noteService, &app, helper)
	userHttp := handlers.NewUserHttpHandler(userService, &app)

	newServer := router.Routes()

	handlers.SetupRoutes(newServer, userHttp, noteHttp)

	err = http.ListenAndServe(":8000", newServer)
	if err != nil {
		log.Fatal("Start server Error")
	}

}
