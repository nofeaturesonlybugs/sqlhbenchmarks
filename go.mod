module github.com/nofeaturesonlybugs/sqlhbenchmarks

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Masterminds/squirrel v1.5.0
	github.com/georgysavva/scany v0.2.8
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.2
	github.com/nofeaturesonlybugs/errors v1.0.1
	github.com/nofeaturesonlybugs/set v0.3.0
	github.com/nofeaturesonlybugs/sqlh v0.1.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
	modernc.org/sqlite v1.10.8
)

replace github.com/nofeaturesonlybugs/sqlh v0.1.0 => ../sqlh
