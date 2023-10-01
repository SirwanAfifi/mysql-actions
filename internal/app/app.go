package app

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/sirwanafifi/mysql-actions/internal/config"
)

func ExecuteActions(eventChan <-chan string, config config.Config, db *sql.DB) {
	for event := range eventChan {
		log.Printf("New event: %s", event)
		for _, job := range config.Jobs {
			log.Printf("Executing job: %s", job.Name)
			steps := job.Steps
			for _, step := range steps {
				log.Printf("Executing step: %s", step.Name)
				shell := step.Shell

				cmd := exec.Command(shell, "-c", step.Run)
				out, err := cmd.Output()
				if err != nil {
					fmt.Println("could not run command: ", err)
				}
				fmt.Println("Output: ", string(out))
			}
		}
	}
}

func PollEventLog(db *sql.DB, config config.Config, eventChan chan<- string) {
	query := "SELECT id, event_type, table_name FROM event_log WHERE id > ? ORDER BY id ASC"

	var lastEventID int64
	err := db.QueryRow("SELECT IFNULL(MAX(id), 0) FROM event_log").Scan(&lastEventID)
	if err != nil {
		log.Fatalf("error fetching the current maximum event ID: %v", err)
	}

	for {
		rows, err := db.Query(query, lastEventID)
		if err != nil {
			log.Printf("error fetching events: %v", err)
			continue
		}

		var newLastEventID int64
		newEvents := false
		for rows.Next() {
			newEvents = true
			var eventID int64
			var eventType, tableName string

			if err := rows.Scan(&eventID, &eventType, &tableName); err != nil {
				log.Printf("error scanning event row: %v", err)
				continue
			}

			eventChan <- fmt.Sprintf("%s on %s", eventType, tableName)
			newLastEventID = eventID
		}

		if newEvents {
			lastEventID = newLastEventID
		}

		rows.Close()
		time.Sleep(1 * time.Second)
	}
}

func CreatedEventLogTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS event_log (
		id bigint unsigned NOT NULL AUTO_INCREMENT,
		event_type enum('insert','update','delete') NOT NULL,
		table_name varchar(255) NOT NULL,
		created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	  )
	`)
	if err != nil {
		return fmt.Errorf("error creating event log table: %v", err)
	}

	return nil
}

func CreateTriggers(db *sql.DB, config config.Config) error {
	triggerTemplate := "CREATE TRIGGER IF NOT EXISTS `%s` " +
		"AFTER %s ON `%s`" +
		"FOR EACH ROW " +
		"INSERT INTO event_log (event_type, table_name) " +
		"VALUES ('%s', '%s');"

	for event, eventConfig := range config.On {
		for _, tableName := range eventConfig.Tables {
			triggerName := fmt.Sprintf("%s_%s_trigger", event, tableName)

			_, err := db.Exec(fmt.Sprintf(triggerTemplate, triggerName, event, tableName, event, tableName))
			if err != nil {
				return fmt.Errorf("error creating trigger for %s on %s: %v", event, tableName, err)
			}
		}
	}

	return nil
}
