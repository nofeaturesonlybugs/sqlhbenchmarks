package sqlhbenchmarks_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/nofeaturesonlybugs/sqlh/grammar"
	"github.com/nofeaturesonlybugs/sqlhbenchmarks"
)

func BenchmarkLibpqSelect(b *testing.B) {
	skip, db, gb, err := sqlhbenchmarks.ConnectLibpq()
	if skip != "" {
		b.Skipf("skipping -- " + skip)
	} else if err != nil {
		b.Fatalf("connect to Postgres failed with %v", err.Error())
	}
	//
	addresses, mdb, err := sqlhbenchmarks.LibpqModels()
	if err != nil {
		b.Fatalf("getting models failed with %v", err.Error())
	}
	if err = mdb.Insert(db, addresses); err != nil {
		b.Fatalf("seeding database with %v", err.Error())
	}
	//
	limits := []int{
		5,
		50,
		100,
		500,
		1000,
	}
	for _, limit := range limits {
		b.Run(fmt.Sprintf("database/sql %v rows", limit), sqlhbenchmarks.StandardSelect(limit, db))
		b.Run(fmt.Sprintf("GORM %v rows", limit), sqlhbenchmarks.GORMSelect(limit, gb))
		b.Run(fmt.Sprintf("sqlx %v rows", limit), sqlhbenchmarks.SqlxSelect(limit, db))
		b.Run(fmt.Sprintf("scany %v rows", limit), sqlhbenchmarks.ScanySelect(limit, db))
		b.Run(fmt.Sprintf("sqlh %v rows", limit), sqlhbenchmarks.SqlhSelect(limit, db))
	}
}

func BenchmarkLibpqInsert(b *testing.B) {
	skip, db, gb, err := sqlhbenchmarks.ConnectLibpq()
	if skip != "" {
		b.Skipf("skipping -- " + skip)
	} else if err != nil {
		b.Fatalf("connect to Postgres failed with %v", err.Error())
	}
	//
	addresses, mdb, err := sqlhbenchmarks.LibpqModels()
	if err != nil {
		b.Fatalf("getting models failed with %v", err.Error())
	}
	//
	limits := []int{
		5,
		50,
		100,
		500,
		1000,
	}
	//
	b.ResetTimer()
	for _, lim := range limits {
		b.Run(fmt.Sprintf("database/sql insert %v row(s)", lim), sqlhbenchmarks.StandardInsert(addresses[0:lim], grammar.Postgres, db))
		b.Run(fmt.Sprintf("GORM insert %v row(s)", lim), sqlhbenchmarks.GORMInsert(addresses[0:lim], gb))
		b.Run(fmt.Sprintf("squirrel insert %v row(s)", lim), sqlhbenchmarks.SquirrelInsert(addresses[0:lim], db))
		b.Run(fmt.Sprintf("sqlh/model insert %v row(s)", lim), sqlhbenchmarks.ModelInsert(mdb, addresses[0:lim], db))
	}
}

func BenchmarkLibpqPreparedInsert(b *testing.B) {
	skip, db, gb, err := sqlhbenchmarks.ConnectLibpq()
	if skip != "" {
		b.Skipf("skipping -- " + skip)
	} else if err != nil {
		b.Fatalf("connect to Postgres failed with %v", err.Error())
	}
	//
	addresses, mdb, err := sqlhbenchmarks.LibpqModels()
	if err != nil {
		b.Fatalf("getting models failed with %v", err.Error())
	}
	//
	limits := []int{
		5,
		50,
		100,
		500,
		1000,
	}
	//
	b.ResetTimer()
	for _, lim := range limits {
		b.Run(fmt.Sprintf("database/sql begin+prepare+insert %v row(s)", lim), sqlhbenchmarks.StandardPreparedInsert(addresses[0:lim], grammar.Postgres, db))
		b.Run(fmt.Sprintf("GORM slice+insert %v row(s)", lim), sqlhbenchmarks.GORMPreparedInsert(addresses[0:lim], gb))
		b.Run(fmt.Sprintf("squirrel begin+prepare+insert %v row(s)", lim), sqlhbenchmarks.SquirrelPreparedInsert(addresses[0:lim], db))
		b.Run(fmt.Sprintf("sqlh/model begin+prepare+insert %v row(s)", lim), sqlhbenchmarks.ModelPreparedInsert(mdb, addresses[0:lim], db))
	}
}

func BenchmarkLibpqUpdate(b *testing.B) {
	skip, db, gb, err := sqlhbenchmarks.ConnectLibpq()
	if skip != "" {
		b.Skipf("skipping -- " + skip)
	} else if err != nil {
		b.Fatalf("connect to Postgres failed with %v", err.Error())
	}
	//
	addresses, mdb, err := sqlhbenchmarks.LibpqModels()
	if err != nil {
		b.Fatalf("getting models failed with %v", err.Error())
	}
	//
	if err = mdb.Insert(db, addresses); err != nil {
		b.Fatalf("seeding database with %v", err.Error())
	}
	// Now modify every address.
	for _, address := range addresses {
		address.Street = address.Street + address.Street
		address.City = address.City + address.City
		address.State = address.State + address.State
		address.Zip = address.Zip + address.Zip
		address.ModifiedTime.Time = address.ModifiedTime.Time.Add(-1 * time.Hour) // Just to make sure modified time updates
	}
	//
	limits := []int{
		5,
		50,
		100,
		500,
		1000,
	}
	b.ResetTimer()
	var tx *sql.Tx
	for _, lim := range limits {
		if tx, err = db.Begin(); err != nil {
			b.Fatalf("pg failed with begin %v", err.Error())
		}
		b.Run(fmt.Sprintf("database/sql update %v row(s)", lim), sqlhbenchmarks.StandardUpdate(addresses[0:lim], grammar.Postgres, tx))
		if err = tx.Rollback(); err != nil {
			b.Fatalf("pg failed with rollback %v", err.Error())
		}
		//
		{
			tx := gb.Begin()
			b.Run(fmt.Sprintf("GORM update %v row(s)", lim), sqlhbenchmarks.GORMUpdate(addresses[0:lim], tx))
			tx.Rollback()
		}
		//
		if tx, err = db.Begin(); err != nil {
			b.Fatalf("pg failed with begin %v", err.Error())
		}
		b.Run(fmt.Sprintf("squirrel update %v row(s)", lim), sqlhbenchmarks.SquirrelUpdate(addresses[0:lim], tx))
		if err = tx.Rollback(); err != nil {
			b.Fatalf("pg failed with rollback %v", err.Error())
		}
		//
		if tx, err = db.Begin(); err != nil {
			b.Fatalf("pg failed with begin %v", err.Error())
		}
		b.Run(fmt.Sprintf("sqlh/model update %v row(s)", lim), sqlhbenchmarks.ModelUpdate(mdb, addresses[0:lim], tx))
		if err = tx.Rollback(); err != nil {
			b.Fatalf("pg failed with rollback %v", err.Error())
		}
	}
}

func BenchmarkLibpqPreparedUpdate(b *testing.B) {
	skip, db, gb, err := sqlhbenchmarks.ConnectLibpq()
	if skip != "" {
		b.Skipf("skipping -- " + skip)
	} else if err != nil {
		b.Fatalf("connect to Postgres failed with %v", err.Error())
	}
	//
	addresses, mdb, err := sqlhbenchmarks.LibpqModels()
	if err != nil {
		b.Fatalf("getting models failed with %v", err.Error())
	}
	//
	if err = mdb.Insert(db, addresses); err != nil {
		b.Fatalf("seeding database with %v", err.Error())
	}
	// Now modify every address.
	for _, address := range addresses {
		address.Street = address.Street + address.Street
		address.City = address.City + address.City
		address.State = address.State + address.State
		address.Zip = address.Zip + address.Zip
		address.ModifiedTime.Time = address.ModifiedTime.Time.Add(-1 * time.Hour) // Just to make sure modified time updates
	}
	//
	limits := []int{
		5,
		50,
		100,
		500,
		1000,
	}
	b.ResetTimer()
	var tx *sql.Tx
	for _, lim := range limits {
		if tx, err = db.Begin(); err != nil {
			b.Fatalf("pg failed with begin %v", err.Error())
		}
		b.Run(fmt.Sprintf("database/sql begin+prepare+update %v row(s)", lim), sqlhbenchmarks.StandardPreparedUpdate(addresses[0:lim], grammar.Postgres, tx))
		if err = tx.Rollback(); err != nil {
			b.Fatalf("pg failed with rollback %v", err.Error())
		}
		//
		{
			tx := gb.Begin()
			b.Run(fmt.Sprintf("GORM update %v row(s)", lim), sqlhbenchmarks.GORMPreparedUpdate(addresses[0:lim], tx))
			tx.Rollback()
		}
		//
		if tx, err = db.Begin(); err != nil {
			b.Fatalf("pg failed with begin %v", err.Error())
		}
		b.Run(fmt.Sprintf("squirrel begin+prepare+update %v row(s)", lim), sqlhbenchmarks.SquirrelPreparedUpdate(addresses[0:lim], tx))
		if err = tx.Rollback(); err != nil {
			b.Fatalf("pg failed with rollback %v", err.Error())
		}
		//
		if tx, err = db.Begin(); err != nil {
			b.Fatalf("pg failed with begin %v", err.Error())
		}
		b.Run(fmt.Sprintf("sqlh/model begin+prepare+update %v row(s)", lim), sqlhbenchmarks.ModelPreparedUpdate(mdb, addresses[0:lim], tx))
		if err = tx.Rollback(); err != nil {
			b.Fatalf("pg failed with rollback %v", err.Error())
		}
		b.StopTimer()
		for _, address := range addresses[0:lim] {
			// Due to how model package works we need to reset the modify times here.
			address.ModifiedTime.Time = address.ModifiedTime.Time.Add(-1 * time.Hour)
		}
		b.StartTimer()
	}
}
