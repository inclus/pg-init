package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/avast/retry-go"
	"github.com/spf13/viper"
)

//BuildConnectionString builds the postgresql connection string from viper configuration
func BuildConnectionString() string {
	databaseHostname := viper.GetString("db.host")
	databasePort := viper.GetString("db.port")
	databaseName := viper.GetString("db.database")
	databaseUser := viper.GetString("db.user")
	databasePassword := viper.GetString("db.password")
	extraDatabaseConfig := viper.GetString("db.extra")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s", databaseUser, databasePassword, databaseHostname, databasePort, databaseName, extraDatabaseConfig)

}

//Handle interface that describes the methods on our database dependency
//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6  . Handle
type Handle interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

//WaitForWorkingConnection retries till you can make a connection to the database
func WaitForWorkingConnection(db Handle) error {
	return retry.Do(
		func() error {
			_, err := db.Query("SELECT 1")
			if err != nil {
				return fmt.Errorf("failed to check connection %w", err)
			}
			return nil
		},
		retry.Attempts(viper.GetUint("retry.attempts")),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("#%d: %s\n", n, err)
		}),
	)
}

//SetupPostgisExtension sets up the postgis extension if not present
func SetupPostgisExtension(db Handle) error {
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS postgis")
	return err
}
