package backend

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq"

	"secondHand/constants"
)

var (
	DBBackend *PostgreSQLBackend
)

type PostgreSQLBackend struct {
	db *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = constants.POSTGRES_PASSWORD
	dbname   = constants.POSTGRES_DB
)

func InitPostgreSQLBackend() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	// check db
	if db.Ping() != nil {
		panic(err)
	}

	fmt.Println("Connected to PostgreSQL database successfully!")

	sqlStrings, err := ioutil.ReadFile("../resources/database-init.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(sqlStrings))
	if err != nil {
		panic(err)
	}

	DBBackend = &PostgreSQLBackend{db: db}
	fmt.Println("Initialized PostgreSQL database successfully!")
}
