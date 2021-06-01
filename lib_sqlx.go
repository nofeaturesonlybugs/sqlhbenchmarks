package sqlhbenchmarks

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/nofeaturesonlybugs/sqlhbenchmarks/types"
)

// SqlxSelectSqlmock creates a test for selecting and scanning rows with sqlx.
func SqlxSelectSqlmock(limit int, mock sqlmock.Sqlmock, db *sql.DB) func(*testing.B) {
	fn := func(b *testing.B) {
		var err error
		var dest []*types.SaleReport
		dbx := sqlx.NewDb(db, "postgres")
		//
		b.StopTimer()
		mockrows := (&types.SaleReport{}).MockRows(limit)
		b.StartTimer()
		//
		for k := 0; k < b.N; k++ {
			b.StopTimer()
			mock.ExpectQuery("select +").WillReturnRows(mockrows)
			dest = nil // Reset dest
			b.StartTimer()
			//
			err = dbx.Select(&dest, "select * from table")
			if err != nil {
				b.Fatalf("sqlx select failed with %v", err.Error())
			}
		}
	}
	return fn
}

// SqlxSelect creates a test for selecting and scanning rows with sqlx.
func SqlxSelect(limit int, db *sql.DB) func(*testing.B) {
	fn := func(b *testing.B) {
		var err error
		var dest []*types.Address
		dbx := sqlx.NewDb(db, "postgres")
		//
		query := `
			select
				pk, created_tmz, modified_tmz,
				street, city, state, zip
			from %v
			limit %v
		`
		query = fmt.Sprintf(query, types.AddressTableName, limit)
		//
		for k := 0; k < b.N; k++ {
			err = dbx.Select(&dest, query)
			if err != nil {
				b.Fatalf("sqlx select failed with %v", err.Error())
			}
		}
	}
	return fn
}
