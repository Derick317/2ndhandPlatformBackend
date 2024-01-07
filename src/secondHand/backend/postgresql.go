package backend

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"secondHand/constants"
	"secondHand/model"
	"secondHand/util"

	// Note: If connecting using the App Engine Flex Go runtime, use
	// "github.com/jackc/pgx/stdlib" instead, since v5 requires
	// Go modules which are not supported by App Engine Flex.
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	dbBackend *PostgreSQLBackend
)

type PostgreSQLBackend struct {
	db *gorm.DB
}

const (
	port = 5432
)

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(db *sql.DB) {
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(7)
	db.SetConnMaxLifetime(1800 * time.Second)
}

// connectUnixSocket initializes a Unix socket connection pool for
// a Cloud SQL instance of Postgres.
func connectUnixSocket() (*sql.DB, error) {
	// Note: Saving credentials in environment variables is convenient, but not
	// secure - consider a more secure solution such as
	// Cloud Secret Manager (https://cloud.google.com/secret-manager) to help
	// keep secrets safe.
	var (
		dbUser         = util.MustGetenv("DB_USER")              // e.g. 'my-db-user'
		dbPwd          = util.MustGetenv("DB_PASS")              // e.g. 'my-db-password'
		unixSocketPath = util.MustGetenv("INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
		dbName         = util.MustGetenv("DB_NAME")              // e.g. 'my-database'
	)

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		dbUser, dbPwd, dbName, unixSocketPath)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	configureConnectionPool(dbPool)

	return dbPool, nil
}

func connectToPostgreSQL() (*gorm.DB, error) {
	if constants.DEPLOYED {
		sqlDB, err := connectUnixSocket()
		if err != nil {
			log.Fatalf("connectUnixSocket: unable to connect: %s", err)
		}
		return gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
	}

	if constants.LOCAL_POSTGRES {
		psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			"localhost", port, util.MustGetenv("LOCAL_POSTGRES_USER"),
			util.MustGetenv("LOCAL_POSTGRES_PASSWORD"), util.MustGetenv("LOCAL_POSTGRES_DB"))
		return gorm.Open(postgres.Open(psqlconn), &gorm.Config{})
	}
	// connect to google
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		util.MustGetenv("GOOGLE_INSTANCE_CONNECTION_NAME"), util.MustGetenv("GOOGLE_POSTGRES_USER"),
		util.MustGetenv("GOOGLE_POSTGRES_DB"), util.MustGetenv("GOOGLE_POSTGRES_PASSWORD"))
	return gorm.Open(postgres.New(postgres.Config{
		DriverName: "cloudsqlpostgres",
		DSN:        dsn,
	}), &gorm.Config{Logger: logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{SlowThreshold: time.Second},
	),
	})
}

func addInitRecords(db *gorm.DB) error {
	var err error = nil
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
	if err = db.Migrator().DropTable(&model.Order{}); err != nil {
		return err
	}
	if err = db.AutoMigrate(&model.Order{}); err != nil {
		return err
	}

	var user = model.User{Email: "alice@alice", Username: "Alice", Password: "alice"}
	db.Create(&user)
	db.Create(&model.User{Email: "bob@bob.com", Username: "Bob", Password: "bob"})
	fmt.Println("Initialized PostgreSQL database successfully!")
	return err
}

func InitPostgreSQLBackend() error {
	db, err := connectToPostgreSQL()
	if err != nil {
		return err
	}

	fmt.Println("Connected to PostgreSQL database successfully!")
	// if err = addInitRecords(db); err != nil {
	// 	return err
	// }
	dbBackend = &PostgreSQLBackend{db: db}
	return nil
}

// CreateRecord creates a user or item or order record in the corresponding table,
// and returns the inserted data's primary key in value's id.
// VALUE should be a pointer to an instance of user or item model.
func CreateRecord(value interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = dbBackend.db
	}
	return tx.Create(value).Error
}

// The record will be save in DEST, so DEST should be a pointer.
func ReadFromDBByPrimaryKey(dest interface{}, primaryKey interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = dbBackend.db
	}
	result := tx.First(dest, primaryKey)
	return result.Error
}

// ReadFromDBByKey queries the database by a key. The record will be save in DEST, so
// DEST should be a pointer, no matter whether its a structure or a slice.
// It returns gorm.ErrRecordNotFound if no record is found.
func ReadFromDBByKey(dest interface{}, key string, target interface{}, onlyFirst bool,
	tx *gorm.DB) error {
	if tx == nil {
		tx = dbBackend.db
	}
	var result = tx.Where(fmt.Sprintf("%s = ?", key), target)
	if onlyFirst {
		result = result.First(dest)
	} else {
		result = result.Find(dest)
	}
	return result.Error
}

// ReadFromDBByKeys queries the database by keys. The record will be save in DEST, so
// DEST should be a pointer, no matter whether its a structure or a slice.
// It returns gorm.ErrRecordNotFound if no record is found.
func ReadFromDBByKeys(dest interface{}, keys []string, targets []string, onlyFirst bool,
	tx *gorm.DB) error {
	if tx == nil {
		tx = dbBackend.db
	}
	var result = tx
	for idx, key := range keys {
		result = result.Where(fmt.Sprintf("%s = ?", key), targets[idx])
	}
	if onlyFirst {
		result = result.First(dest)
	} else {
		result = result.Find(dest)
	}
	return result.Error
}

// ReadFromDBEqualOrInclude queries record(s) whose KEY(s) equal or include corresponding
// TARGETS (non-case sensitivity). The record will be save in DEST, so
// DEST should be a pointer, no matter whether its a structure or a slice.
// It returns gorm.ErrRecordNotFound if no record is found.
func ReadFromDBEqualOrInclude(dest interface{}, keys []string, targets []interface{},
	equal []bool, onlyFirst bool, tx *gorm.DB) error {
	if tx == nil {
		tx = dbBackend.db
	}
	var result = tx
	for idx, key := range keys {
		if equal[idx] {
			result = result.Where(fmt.Sprintf("%s = ?", key), targets[idx])
		} else {
			for _, target := range (targets[idx]).([]string) {
				result = result.Where(fmt.Sprintf("%s iLike ?", key),
					fmt.Sprintf("%%%s%%", target))
			}
		}
	}
	if onlyFirst {
		result = result.First(dest)
	} else {
		result = result.Find(dest)
	}
	return result.Error
}

// DEST should be a pointer. It has no function but to represent table.
func DeleteFromDBByPrimaryKey(dest interface{}, primaryKey interface{}, tx *gorm.DB) error {
	if tx == nil {
		tx = dbBackend.db
	}
	result := tx.Delete(dest, primaryKey)
	return result.Error
}

// If the specified value has no primary value, a batch delete will be performed.
// It will delete all matched records.
// DEST should be a pointer. If DEST contains primary key, it is included in the conditions.
//
// It reports the number of deleted records and possible errors
func DeleteFromDBByKey(dest interface{}, key string, target interface{},
	tx *gorm.DB) (int64, error) {
	if tx == nil {
		tx = dbBackend.db
	}
	var result = tx.Where(fmt.Sprintf("%s = ?", key), target).Delete(dest)
	return result.RowsAffected, result.Error
}

// If the specified value has no primary value, a batch delete will be performed.
// It will delete all matched records.
// DEST should be a pointer. If DEST contains primary key, it is included in the conditions.
//
// It reports the number of deleted records and possible errors
func DeleteFromDBByKeys(dest interface{}, keys []string, targets []string,
	tx *gorm.DB) (int64, error) {
	if tx == nil {
		tx = dbBackend.db
	}
	var result = tx
	for idx, key := range keys {
		result = result.Where(fmt.Sprintf("%s = ?", key), targets[idx])
	}
	result = result.Delete(dest)
	return result.RowsAffected, result.Error
}

// UpdateColumnsWithConditions updates columns of record(s).
// We update all records which match dest's non-zero field and the condition (key and target).
// DEST should be a pointer.
// If VALUE is nil, we just save DEST as it is.
// It returns the number of affected rows and possible error.
func UpdateColumnsWithConditions(dest interface{}, key string, target interface{},
	value map[string]interface{}, tx *gorm.DB) (int64, error) {
	if tx == nil {
		tx = dbBackend.db
	}
	var result *gorm.DB
	if value == nil {
		result = tx.Save(dest)
	} else if key == "" || target == nil {
		result = tx.Model(dest).Updates(value)
	} else {
		result = tx.Model(dest).Where(key+" = ?", target).Updates(value)
	}
	return result.RowsAffected, result.Error
}

// BeginTransaction begins a transaction.
func BeginTransaction() *gorm.DB {
	return dbBackend.db.Begin()
}
