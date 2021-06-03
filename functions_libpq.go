package sqlhbenchmarks

import (
	"database/sql"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/lib/pq"
	"github.com/nofeaturesonlybugs/errors"
	"github.com/nofeaturesonlybugs/sqlh/grammar"
	"github.com/nofeaturesonlybugs/sqlh/model"
	"github.com/nofeaturesonlybugs/sqlhbenchmarks/types"
)

// ConnectLibpq connects to postgresql using lib/pq if the TEST_POSTGRES environment variable is set.
func ConnectLibpq() (SkipReason string, DB *sql.DB, GB *gorm.DB, err error) {
	env := "TEST_POSTGRES"
	//
	// Extra for GORM.
	var gdb *sql.DB
	gcfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	//
	dsn := os.Getenv(env)
	if dsn == "" {
		SkipReason = env + " environment variable is empty"
	} else if DB, err = sql.Open("postgres", dsn); err != nil {
		return
	} else if GB, err = gorm.Open(postgres.Open(dsn), gcfg); err != nil {
		return
	} else if err = DB.Ping(); err != nil {
		return
	} else if gdb, err = GB.DB(); err != nil {
		return
	} else if err = gdb.Ping(); err != nil {
		return
	} else if err = ExecSchema(SchemaLibpq, DB); err != nil {
		return
	}
	return
}

// LibpqModels returns all the models and types for our tests.
func LibpqModels() (Addresses []*types.Address, Mdb *model.Models, err error) {
	Addresses = types.AddressRecords
	//
	Mdb = types.NewModels(grammar.Postgres)
	if Mdb == nil {
		err = errors.Errorf("nil mdb")
		return
	}
	return
}
