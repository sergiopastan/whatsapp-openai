package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"github.com/sergiopastan/whatsapp-openai/config"
)

func Connect(conf config.DbConfig) (*sql.DB, error) {
	log.Infof("conf: %v", conf)
	db, err := sql.Open("sqlite3", conf.DBHost)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	db.SetMaxOpenConns(conf.DBPoolSize)
	db.SetMaxIdleConns(conf.DBIdlePoolSize)
	db.SetConnMaxLifetime(conf.DBConnLifetime)
	return db, nil
}
