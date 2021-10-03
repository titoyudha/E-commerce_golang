package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
}

type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

func (server *Server) Initialize(appConfig AppConfig, dbCongig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	var err error
	dsn := "user:pass@tcp(127.0.0.1:3306)/go_shop?charset=utf8mb4&parseTime=True&loc=Local"

	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	server.Router = mux.NewRouter()
	server.InitializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Run() {
	var server = Server{}
	var appConfig = AppConfig{}

	var dbCongig = DBConfig{}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	appConfig.AppName = getEnv("APP_NAME", "default app name")
	appConfig.AppEnv = getEnv("APP_ENV", "default app env")
	appConfig.AppPort = getEnv("APP_PORT", "default app port")

	dbCongig.DBHost = getEnv("DB_HOST", "localhost")
	dbCongig.DBUser = getEnv("DB_USER", "go_shop")
	dbCongig.DBPassword = getEnv("DB_PASSWORD", "")
	dbCongig.DBName = getEnv("DB_NAME", "go_shop")
	dbCongig.DBPort = getEnv("DB_PORT", "3306")

	server.Initialize(appConfig, dbCongig)
	server.Run(":" + appConfig.AppPort)
}
