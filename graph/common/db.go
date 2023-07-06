// @/graph/common/db.go
package common

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
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

	db.Use(prometheus.New(prometheus.Config{
		DBName:          "db1", // `DBName` as metrics label
		RefreshInterval: 15,    // refresh metrics interval (default 15 seconds)
		//PushAddr:        "localhost:9100", // push metrics if `PushAddr` configured
		StartServer:    true, // start http server to expose metrics
		HTTPServerPort: 8081, // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
		//MetricsCollector: []prometheus.MetricsCollector{
		//	&prometheus.Postgres{VariableNames: []string{"Threads_running"}},
		//},
	}))

	if err != nil {
		panic("failed to connect database")
	}
	return db, nil
}
