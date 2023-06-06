package api_tests

import (
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/leavemeal0ne/Golang-2023/internal/config"
	"github.com/leavemeal0ne/Golang-2023/internal/driver"
	"github.com/leavemeal0ne/Golang-2023/internal/handlers"
	"github.com/leavemeal0ne/Golang-2023/internal/middleware"
	"github.com/leavemeal0ne/Golang-2023/internal/repositories"
	router2 "github.com/leavemeal0ne/Golang-2023/internal/router"
	"github.com/leavemeal0ne/Golang-2023/internal/services"
	"github.com/steinfletcher/apitest"
	"log"
	"net/http"
	"os"
	"testing"
)

var dbDATA string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file, using system env variables")
	}

	dbDATA = os.Getenv("DATABASE_DATA")
}

func TestGetCreate(t *testing.T) {
	//міграції та очистку бд перед тестами виконую за допомоги команди в make файлі
	DB, err := driver.ConnectSQL(dbDATA)
	if err != nil {
		log.Println(err)
		log.Fatal("Start server error")
	}
	defer func(SQL *sql.DB) {
		err := SQL.Close()
		if err != nil {

		}
	}(DB.SQL)

	app := config.NewAppConfig()
	middleware.SetSession(app.Session)

	userRepo := repositories.NewUserRepository(DB)
	userService := services.NewUserService(userRepo)
	userHttp := handlers.NewUserHttpHandler(userService, app)

	router := router2.Routes()
	handlers.SetupRoutes(router, userHttp, nil)

	apitest.New().
		Handler(router).
		Post("/api/user/signup").
		Body(`{"email": "jorik246@ukr.net","password_hash": "121"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	apitest.New().
		Handler(router).
		Post("/api/user/signup").
		Body(`{"email": "jorik246@ukr.net","password_hash": "120921"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

	apitest.New().
		Handler(router).
		Post("/api/user/signup").
		Body(`{"email": "jorik246@ukr.net","password_hash": "120921"}`).
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	apitest.New().
		Handler(router).
		Post("/api/user/login").
		Body(`{"email": "jorik246@ukr.net","password_hash": "120921"}`).
		Expect(t).
		Status(http.StatusOK).
		End()

}
