package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
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
	DBDriver   string
}

func (server *Server) Initialize(appConfig AppConfig, dbCongig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	server.initializeDB(dbCongig)

	server.InitializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func (server *Server) initializeDB(dbConfig DBConfig) {
	var err error

	if dbConfig.DBDriver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
		server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable Timezone=Asia/Jakarta",
			dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)

		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		panic(err)
	}
	for _, models := range RegisterModels() {
		err = server.DB.Debug().AutoMigrate(models.Model)

		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Database Succesfully Migrated")
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
	dbCongig.DBDriver = getEnv("DB_DRIVER", "postgresql")

	server.Initialize(appConfig, dbCongig)
	server.Run(":" + appConfig.AppPort)
}
