package db

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Database struct {
	Client *sqlx.DB
}

func NewDatabase(config *viper.Viper) *Database {
	// MySQL connection string format
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.host"),
		config.GetInt("database.port"),
		config.GetString("database.name"),
	)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("error when connecting to database: %v", err)
	}

	return &Database{
		Client: db,
	}
}
