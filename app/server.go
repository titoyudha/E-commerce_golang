package app

import (
	"flag"
	"fmt"
	"go_shop/database/seeders"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
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
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		panic("Failed on connecting to the database server")
	}
}

func (server *Server) dbMigrate() {
	for _, models := range RegisterModels() {
		err := server.DB.Debug().AutoMigrate(models.Model)

		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Database Succesfully Migrated")
}

func (server *Server) initCommands(config AppConfig, dbConfig DBConfig) {
	server.initializeDB(dbConfig)

	cmdApp := cli.NewApp()
	cmdApp.Commands = []*cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				server.dbMigrate()
				return nil
			},
		},
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(server.DB)
				if err != nil {
					panic(err)
				}
				return nil
			},
		},
	}
	err := cmdApp.Run(os.Args)
	if err != nil {
		panic(err)
	}
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

	flag.Parse()
	arg := flag.Arg(0)
	if arg != "" {
		server.initCommands(appConfig, dbCongig)
	} else {
		server.Initialize(appConfig, dbCongig)
		server.Run(":" + appConfig.AppPort)
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
