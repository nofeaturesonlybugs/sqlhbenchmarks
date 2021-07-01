package sqlhbenchmarks

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nofeaturesonlybugs/sqlh/grammar"
	"github.com/nofeaturesonlybugs/sqlhbenchmarks/types"
)

// StandardSelectSqlmock creates a test for selecting and scanning rows with database/sql.
func StandardSelectSqlmock(limit int, mock sqlmock.Sqlmock, db *sql.DB) func(*testing.B) {
	fn := func(b *testing.B) {
		var rows *sql.Rows
		var err error
		var d *types.SaleReport
		//
		b.StopTimer()
		mockrows := (&types.SaleReport{}).MockRows(limit)
		b.StartTimer()
		//
		for k := 0; k < b.N; k++ {
			b.StopTimer()
			mock.ExpectQuery("select +").WillReturnRows(mockrows)
			b.StartTimer()
			//
			rows, err = db.Query("select * from table")
			if err != nil {
				b.Fatalf("database/sql query failed with %v", err.Error())
			}
			for rows.Next() {
				d = &types.SaleReport{}
				err = rows.Scan(
					&d.Id, &d.CreatedTime, &d.ModifiedTime,
					&d.Price, &d.Quantity, &d.Total,
					&d.CustomerId, &d.CustomerFirst, &d.CustomerLast,
					&d.VendorId, &d.VendorName, &d.VendorDescription,
					&d.VendorContactId, &d.VendorContactFirst, &d.VendorContactLast,
				)
			}
			if err != nil {
				b.Fatalf("database/sql scan failed with %v", err.Error())
			}
			if err = rows.Err(); err != nil {
				b.Fatalf("database/sql rows.Err failed with %v", err.Error())
			}
			rows.Close()
		}
	}
	return fn
}

// StandardSelect creates a test for selecting and scanning rows with database/sql.
func StandardSelect(limit int, db *sql.DB) func(*testing.B) {
	fn := func(b *testing.B) {
		var rows *sql.Rows
		var err error
		var d *types.Address
		//
		query := `
			select
				pk, created_tmz, modified_tmz,
				street, city, state, zip
			from %v
			limit %v
		`
		query = fmt.Sprintf(query, types.AddressTableName, limit)
		for k := 0; k < b.N; k++ {
			rows, err = db.Query(query)
			if err != nil {
				b.Fatalf("database/sql query failed with %v", err.Error())
			}
			for rows.Next() {
				d = &types.Address{}
				err = rows.Scan(
					&d.Id, &d.CreatedTime, &d.ModifiedTime,
					&d.Street, &d.City, &d.State, &d.Zip,
				)
			}
			if err != nil {
				b.Fatalf("database/sql scan failed with %v", err.Error())
			}
			if err = rows.Err(); err != nil {
				b.Fatalf("database/sql rows.Err failed with %v", err.Error())
			}
			rows.Close()
		}
	}
	return fn
}

// StandardInsert performs INSERTs using QueryRow() -> row.Scan() over the range of models using
// standard database/sql package.
func StandardInsert(addresses []*types.Address, g grammar.Grammar, db *sql.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var query string
		switch g {
		case grammar.Sqlite:
			query = `
				insert into %v ( street, city, state, zip )
				values ( ?, ?, ?, ? )
				returning pk, created_tmz, modified_tmz
			`
		case grammar.Postgres:
			query = `
				insert into %v ( street, city, state, zip )
				values ( $1, $2, $3, $4 )
				returning pk, created_tmz, modified_tmz
			`
		}
		query = fmt.Sprintf(query, types.AddressTableName)
		//
		var row *sql.Row
		var err error
		//
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreInsert(b)
				//
				row = db.QueryRow(query, address.Street, address.City, address.State, address.Zip)
				if err = row.Scan(&address.Id, &address.CreatedTime, &address.ModifiedTime); err != nil {
					b.Fatalf("standard failed with %v", err.Error())
				}
				//
				address.PostInsert(b)
			}
		}
	}
	return fn
}

// StandardPreparedInsert performs INSERTs using Prepare() -> QueryRow() -> row.Scan() over the
// range of models using standard database/sql package.
func StandardPreparedInsert(addresses []*types.Address, g grammar.Grammar, db *sql.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var query string
		switch g {
		case grammar.Sqlite:
			query = `
				insert into %v ( street, city, state, zip )
				values ( ?, ?, ?, ? )
				returning pk, created_tmz, modified_tmz
			`
		case grammar.Postgres:
			query = `
				insert into %v ( street, city, state, zip )
				values ( $1, $2, $3, $4 )
				returning pk, created_tmz, modified_tmz
			`
		}
		query = fmt.Sprintf(query, types.AddressTableName)
		//
		var tx *sql.Tx
		var stmt *sql.Stmt
		var row *sql.Row
		var err error
		//
		if tx, err = db.Begin(); err != nil {
			b.Fatalf("error beginning transaction with %v", err.Error())
		}
		defer tx.Rollback()
		if stmt, err = tx.Prepare(query); err != nil {
			b.Fatalf("error preparing statement with %v", err.Error())
		}
		defer stmt.Close()
		//
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				row = stmt.QueryRow(address.Street, address.City, address.State, address.Zip)
				if err = row.Scan(&address.Id, &address.CreatedTime, &address.ModifiedTime); err != nil {
					b.Fatalf("standard failed with %v", err.Error())
				}
			}
		}
		//
		if tx != nil {
			if err = tx.Commit(); err != nil {
				b.Fatalf("error durring commit with %v", err.Error())
			}
		}
	}
	return fn
}

// StandardUpdate performs UPDATEs using QueryRow() -> row.Scan() over the range of models using
// standard database/sql package.
func StandardUpdate(addresses []*types.Address, g grammar.Grammar, tx *sql.Tx) func(b *testing.B) {
	fn := func(b *testing.B) {
		var query string
		switch g {
		case grammar.Sqlite:
			query = `
				update %v set
					street = ?, city = ?, state = ?, zip = ?
				where pk = ?
				returning modified_tmz
			`
		case grammar.Postgres:
			query = `
				update %v set
					street = $1, city = $2, state = $3, zip = $4
				where pk = $5
				returning modified_tmz
			`
		}
		query = fmt.Sprintf(query, types.AddressTableName)
		//
		var row *sql.Row
		var err error
		//
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreUpdate(b)
				//
				row = tx.QueryRow(query, address.Street, address.City, address.State, address.Zip, address.Id)
				if err = row.Scan(&address.ModifiedTime); err != nil {
					b.Fatalf("standard failed with %v", err.Error())
				}
				//
				address.PostUpdate(b)
			}
		}
	}
	return fn
}

// StandardPreparedUpdate performs UPDATEs using QueryRow() -> row.Scan() over the range of models using
// standard database/sql package.
func StandardPreparedUpdate(addresses []*types.Address, g grammar.Grammar, tx *sql.Tx) func(b *testing.B) {
	fn := func(b *testing.B) {
		var query string
		switch g {
		case grammar.Sqlite:
			query = `
				update %v set
					street = ?, city = ?, state = ?, zip = ?
				where pk = ?
				returning modified_tmz
			`
		case grammar.Postgres:
			query = `
				update %v set
					street = $1, city = $2, state = $3, zip = $4
				where pk = $5
				returning modified_tmz
			`
		}
		query = fmt.Sprintf(query, types.AddressTableName)
		//
		var stmt *sql.Stmt
		var row *sql.Row
		var err error
		//
		if stmt, err = tx.Prepare(query); err != nil {
			b.Fatalf("error preparing statement with %v", err.Error())
		}
		defer stmt.Close()

		//
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreUpdate(b)
				//
				row = stmt.QueryRow(address.Street, address.City, address.State, address.Zip, address.Id)
				if err = row.Scan(&address.ModifiedTime); err != nil {
					b.Fatalf("standard failed with %v", err.Error())
				}
				//
				address.PostUpdate(b)
			}
		}
	}
	return fn
}
