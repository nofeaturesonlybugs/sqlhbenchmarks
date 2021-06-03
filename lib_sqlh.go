package sqlhbenchmarks

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nofeaturesonlybugs/sqlh"
	"github.com/nofeaturesonlybugs/sqlh/model"
	"github.com/nofeaturesonlybugs/sqlhbenchmarks/types"
)

// SqlhSelectSqlmock creates a test for selecting and scanning rows with sqlh.
func SqlhSelectSqlmock(limit int, mock sqlmock.Sqlmock, db *sql.DB) func(*testing.B) {
	fn := func(b *testing.B) {
		var err error
		var dest []*types.SaleReport
		scanner := &sqlh.Scanner{
			Mapper: types.NewMapper(),
		}
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
			err = scanner.Select(db, &dest, "select * from table")
			if err != nil {
				b.Fatalf("sqlh select failed with %v", err.Error())
			}
		}
	}
	return fn
}

// SqlhSelect creates a test for selecting and scanning rows with sqlh.
func SqlhSelect(limit int, db *sql.DB) func(*testing.B) {
	fn := func(b *testing.B) {
		var err error
		var dest []*types.Address
		scanner := &sqlh.Scanner{
			Mapper: types.NewMapper(),
		}
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
			err = scanner.Select(db, &dest, query)
			if err != nil {
				b.Fatalf("sqlh select failed with %v", err.Error())
			}
		}
	}
	return fn
}

// ModelInsert performs INSERTs using github.com/nofeaturesonlybugs/sqlh/models package.
func ModelInsert(mdb *model.Models, addresses []*types.Address, db *sql.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var err error
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreInsert(b)
				//
				err = mdb.Insert(db, address)
				if err != nil {
					b.Fatalf("sqlh failed with %v", err.Error())
				}

				//
				address.PostInsert(b)
			}
		}
	}
	return fn
}

// ModelPreparedInsert performs INSERTs using github.com/nofeaturesonlybugs/sqlh/models package by inserting
// the slice, which internally should use a prepared statement.
func ModelPreparedInsert(mdb *model.Models, addresses []*types.Address, db *sql.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var err error
		for k := 0; k < b.N; k++ {
			err = mdb.Insert(db, addresses)
			if err != nil {
				b.Fatalf("sqlh failed with %v", err.Error())
			}
		}
	}
	return fn
}

// ModelUpdate performs UPDATEs using github.com/nofeaturesonlybugs/sqlh/models package.
func ModelUpdate(mdb *model.Models, address []*types.Address, tx *sql.Tx) func(b *testing.B) {
	fn := func(b *testing.B) {
		var err error
		for k := 0; k < b.N; k++ {
			for _, address := range address {
				address.PreUpdate(b)
				//
				err = mdb.Update(tx, address)
				if err != nil {
					b.Fatalf("sqlh failed with %v", err.Error())
				}
				//
				address.PostUpdate(b)
			}
		}
	}
	return fn
}

// ModelPreparedUpdate performs UPDATESs using github.com/nofeaturesonlybugs/sqlh/models package by inserting
// the slice, which internally should use a prepared statement.
func ModelPreparedUpdate(mdb *model.Models, addresses []*types.Address, tx *sql.Tx) func(b *testing.B) {
	fn := func(b *testing.B) {
		var err error
		for k := 0; k < b.N; k++ {
			err = mdb.Update(tx, addresses)
			if err != nil {
				b.Error(err.Error())
				b.FailNow()
			}
		}
	}
	return fn
}
