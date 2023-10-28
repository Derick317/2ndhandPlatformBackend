package backend

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"secondHand/constants"
	"secondHand/model"
)

var (
	dbBackend *PostgreSQLBackend
)

type PostgreSQLBackend struct {
	db *gorm.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = constants.POSTGRES_PASSWORD
	dbname   = constants.POSTGRES_DB
)

func InitPostgreSQLBackend() error {

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(psqlconn), &gorm.Config{})

	if err != nil {
		return err
	}

	fmt.Println("Connected to PostgreSQL database successfully!")

	if err = db.Migrator().DropTable(&model.User{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	if err = db.Migrator().DropTable(&model.Item{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&model.Item{}); err != nil {
		return err
	}

	var user = model.User{Email: "alice@alice", Username: "Alice", Password: "alice"}
	db.Create(&user)

	dbBackend = &PostgreSQLBackend{db: db}
	fmt.Println("Initialized PostgreSQL database successfully!")
	return nil
}

// CreateRecord creates a user or item record in the corresponding table,
// and returns the inserted data's primary key in value's id.
// VALUE should be a pointer to an instance of user or item model.
func CreateRecord(value interface{}) error {
	result := dbBackend.db.Create(value)
	return result.Error
}

func ReadFromDBByPrimaryKey(dest interface{}, primaryKey interface{}) error {
	result := dbBackend.db.Find(&dest, primaryKey)
	return result.Error
}

// ReadFromDBByKey query the database by a key
// It returns gorm.ErrRecordNotFound if no record is found
func ReadFromDBByKey(dest interface{}, key string, target string, onlyFirst bool) error {
	var query string
	if onlyFirst {
		query = fmt.Sprintf("%s = ?", key)
	} else {
		query = fmt.Sprintf("%s <> ?", key)
	}
	result := dbBackend.db.Where(query, target).First(&dest)
	return result.Error
}

// UpdateColumnWithConditions updates a column of record(s).
// We update all records which match dest's non-zero field and the condition (key and target).
// DEST should be a pointer.
// If COLUMN or VALUE is empty, we just save DEST as it is.
// It returns the number of affected rows and possible error.
func UpdateColumnWithConditions(dest interface{}, key string, target interface{},
	column string, value interface{}) (int64, error) {
	var result *gorm.DB
	if column == "" || value == nil {
		result = dbBackend.db.Save(dest)
	} else if key == "" || target == nil {
		result = dbBackend.db.Model(&dest).Update(column, value)
	} else {
		result = dbBackend.db.Model(&dest).Where(key+" = ?", target).Update(column, value)
	}
	return result.RowsAffected, result.Error
}
