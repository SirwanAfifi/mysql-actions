package main

import (
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/sirwanafifi/mysql-actions/internal/app"
	"github.com/sirwanafifi/mysql-actions/internal/config"
	"github.com/sirwanafifi/mysql-actions/internal/db"
)

func main() {

	configFile := "./config/sample.yml"
	config := config.ReadConfigFile(configFile)

	user := goDotEnvVariable("MYSQL_USER")
	password := goDotEnvVariable("MYSQL_PASSWORD")
	host := goDotEnvVariable("MYSQL_HOST")
	dbname := goDotEnvVariable("MYSQL_DB")

	db, err := db.ConnectToMySQL(user, password, host, dbname)
	if err != nil {
		log.Fatalf("error connecting to MySQL: %v", err)
	}
	defer db.Close()

	app.CreateTriggers(db, config)

	eventChan := make(chan string)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		app.PollEventLog(db, config, eventChan)
		wg.Done()
	}()

	go func() {
		app.ExecuteActions(eventChan, config, db)
		wg.Done()
	}()

	wg.Wait()
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
