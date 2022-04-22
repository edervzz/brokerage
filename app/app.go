package app

import (
	"brokerage/domain"
	"brokerage/migrations"
	"brokerage/service"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func Run() {
	client := CreateClientDB()
	router := mux.NewRouter()

	migrateHandler := NewMigrationHandler(migrations.NewMigrationServiceInterface(domain.NewMigration(client)))
	accountHandler := NewAccountHandler(service.NewAccountServiceInterface(domain.NewAccountDB(client)))
	orderHandler := NewOrderHandler(service.NewOrderServiceInterface(domain.NewOrderDB(client)))

	router.HandleFunc("/migration", migrateHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/accounts", accountHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/accounts/{id}/orders", orderHandler.Create).Methods(http.MethodPost)

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8001"
	}

	fmt.Println("listening on :", appPort)
	if err := http.ListenAndServe(":"+appPort, router); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func CreateClientDB() *sqlx.DB {

	server := os.Getenv("DB_SERVER")
	if server == "" {
		server = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("DB_PASS")
	if pass == "" {
		pass = "eder"
	}

	connString := user + ":" + pass + "@tcp(" + server + ":" + port + ")/"
	client, err := sqlx.Open("mysql", connString)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return client
}
