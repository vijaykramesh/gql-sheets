// @/graph/common/db.go
package common

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func InitDb() (*gorm.DB, error) {
	var err error
	// fetch host, user, password, dbname, port from env
	host := "localhost"
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := "9000"
	// connect to db
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db, nil
}
