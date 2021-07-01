package sqlhbenchmarks

import (
	"database/sql"
	"os"

	"github.com/nofeaturesonlybugs/errors"
	"github.com/nofeaturesonlybugs/sqlh/grammar"
	"github.com/nofeaturesonlybugs/sqlh/model"
	"github.com/nofeaturesonlybugs/sqlhbenchmarks/types"

	// Uses modernc library -- don't feel like dealing with cgo.
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
	// "gorm.io/gorm/logger"

	_ "modernc.org/sqlite"
)

// ConnectSqlite connects to sqlite using modernc.org/sqlite if the TEST_SQLITE environment variable is set.
func ConnectSqlite() (SkipReason string, DB *sql.DB /*GB *gorm.DB,*/, err error) {
	env := "TEST_SQLITE"
	//
	// Extra for GORM.
	// var gdb *sql.DB
	// gcfg := &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Silent),
	// }
	//
	dsn := os.Getenv(env)
	if dsn == "" {
		SkipReason = env + " environment variable is empty"
	} else if DB, err = sql.Open("sqlite", dsn); err != nil {
		return
		// } else if GB, err = gorm.Open(sqlite.Open(dsn), gcfg); err != nil {
		// 	return
	} else if err = DB.Ping(); err != nil {
		return
		// } else if gdb, err = GB.DB(); err != nil {
		// 	return
		// } else if err = gdb.Ping(); err != nil {
		// 	return
	} else if err = ExecSchema(SchemaSqlite, DB); err != nil {
		return
	}
	return
}

// SqliteModels returns all the models and types for our tests.
func SqliteModels() (Addresses []*types.Address, Mdb *model.Models, err error) {
	Addresses = types.AddressRecords
	//
	Mdb = types.NewModels(grammar.Sqlite)
	if Mdb == nil {
		err = errors.Errorf("nil mdb")
		return
	}
	return
}
