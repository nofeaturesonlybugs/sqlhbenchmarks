package sqlhbenchmarks_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nofeaturesonlybugs/sqlhbenchmarks"
)

func BenchmarkSqlmockSelect(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("creating sqlmock with %v", err.Error())
	}
	limits := []int{
		5,
		50,
		100,
		500,
		1000,
		10000,
	}
	for _, limit := range limits {
		b.Run(fmt.Sprintf("database/sql %v rows", limit), sqlhbenchmarks.StandardSelectSqlmock(limit, mock, db))
		b.Run(fmt.Sprintf("sqlx %v rows", limit), sqlhbenchmarks.SqlxSelectSqlmock(limit, mock, db))
		b.Run(fmt.Sprintf("scany %v rows", limit), sqlhbenchmarks.ScanySelectSqlmock(limit, mock, db))
		b.Run(fmt.Sprintf("sqlh %v rows", limit), sqlhbenchmarks.SqlhSelectSqlmock(limit, mock, db))
	}
}
