package sqlhbenchmarks

import (
	"database/sql"
	"strings"

	"github.com/nofeaturesonlybugs/sqlhbenchmarks/types"
)

var (
	// SchemaLibpq is the slice of queries needed to init the postgres schema.
	SchemaLibpq = []string{
		`DROP TABLE IF EXISTS {TABLE}`,
		`CREATE TABLE {TABLE} (
			pk serial primary key,
			created_tmz timestamp (6) with time zone not null,
			modified_tmz timestamp (6) with time zone not null,
			street character varying not null,
			city character varying not null,
			state character varying not null,
			zip character varying not null
		)`,
		`create or replace function trg_addresses_insert() returns trigger as $BODY$ 
declare
begin
	new.created_tmz = now();
	new.modified_tmz = now();
	return new;
end;
$BODY$ language plpgsql`,
		`create or replace function trg_addresses_update() returns trigger as $BODY$
declare
begin
	new.modified_tmz = now();
	new.created_tmz = old.created_tmz;
	return new;
end;
$BODY$ language plpgsql`,
		`create trigger trg_addresses_insert before insert on {TABLE}
for each row execute procedure trg_addresses_insert()`,
		`create trigger trg_addresses_update before update on {TABLE}
for each row execute procedure trg_addresses_update()`,
	}
	// SchemaSqlite is the slice of queries needed to init the sqlite schema.
	SchemaSqlite = []string{
		`DROP TABLE IF EXISTS {TABLE}`,
		`CREATE TABLE {TABLE} (
			pk integer primary key,
			created_tmz datetime not null default (strftime('%Y-%m-%d %H:%M:%S', 'now', 'utc')),
			modified_tmz datetime not null default (strftime('%Y-%m-%d %H:%M:%S', 'now', 'utc')),
			street text not null,
			city text not null,
			state text not null,
			zip text not null
		)`,
	}
)

// ExecSchema runs statements to init the schema for tests against db.
func ExecSchema(queries []string, db *sql.DB) error {
	var err error
	for _, query := range queries {
		query = strings.Replace(query, "{TABLE}", types.AddressTableName, -1)
		if _, err = db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
