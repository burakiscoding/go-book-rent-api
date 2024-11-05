package database

import (
	"database/sql"
	"os"
)

const DriverMysql = "mysql"

type Config struct {
	Driver   string
	Username string
	Password string
	Host     string
	Name     string
}

// default database config
func NewConfig() Config {
	return Config{
		Driver:   DriverMysql,
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func NewSQLWithConfig(config Config) (*sql.DB, error) {
	dataSourceName := config.Username + ":" + config.Password + "@(" + config.Host + ")/" + config.Name + "?parseTime=true"

	return sql.Open(config.Driver, dataSourceName)
}

// default sql connection
func NewSQL() (*sql.DB, error) {
	config := NewConfig()
	dataSourceName := config.Username + ":" + config.Password + "@(" + config.Host + ")/" + config.Name + "?parseTime=true"

	return sql.Open(config.Driver, dataSourceName)
}
