package sqlhbenchmarks

import (
	"database/sql"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/nofeaturesonlybugs/sqlhbenchmarks/types"
)

// SquirrelInsert performs INSERTs using github.com/Masterminds/squirrel.
func SquirrelInsert(addresses []*types.Address, db *sql.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var err error
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreInsert(b)
				//
				query := sq.Insert(types.AddressTableName).
					Columns("street", "city", "state", "zip").
					Values(address.Street, address.City, address.State, address.Zip).
					Suffix("RETURNING pk, created_tmz, modified_tmz").
					RunWith(db).
					PlaceholderFormat(sq.Dollar)
				if err = query.QueryRow().Scan(&address.Id, &address.CreatedTime, &address.ModifiedTime); err != nil {
					b.Fatalf("squirrel failed with %v", err.Error())
				}
				//
				address.PostInsert(b)
			}
		}
	}
	return fn
}

// SquirrelPreparedInsert performs INSERTs using github.com/Masterminds/squirrel.
func SquirrelPreparedInsert(addresses []*types.Address, db *sql.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var dbcache *sq.StmtCache
		var tx *sql.Tx
		var err error
		//
		if tx, err = db.Begin(); err != nil {
			b.Fatalf("error beginning transaction with %v", err.Error())
		}
		defer tx.Rollback()
		dbcache = sq.NewStmtCache(tx)
		defer dbcache.Clear()
		//
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				query := sq.Insert(types.AddressTableName).
					Columns("street", "city", "state", "zip").
					Values(address.Street, address.City, address.State, address.Zip).
					Suffix("RETURNING pk, created_tmz, modified_tmz").
					RunWith(dbcache).
					PlaceholderFormat(sq.Dollar)
				if err = query.QueryRow().Scan(&address.Id, &address.CreatedTime, &address.ModifiedTime); err != nil {
					b.Fatalf("squirrel failed with %v", err.Error())
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

// SquirrelUpdate performs UPDATEs using github.com/Masterminds/squirrel.
func SquirrelUpdate(addresses []*types.Address, tx *sql.Tx) func(b *testing.B) {
	fn := func(b *testing.B) {
		var err error
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreUpdate(b)
				//
				query := sq.Update(types.AddressTableName).
					Set("street", address.Street).
					Set("city", address.City).
					Set("state", address.State).
					Set("zip", address.Zip).
					Where(sq.Eq{"pk": address.Id}).
					Suffix("RETURNING modified_tmz").
					RunWith(tx).
					PlaceholderFormat(sq.Dollar)
				if err = query.QueryRow().Scan(&address.ModifiedTime); err != nil {
					b.Fatalf("squirrel failed with %v", err.Error())
				}
				//
				address.PostUpdate(b)
			}
		}
	}
	return fn
}

// SquirrelPreparedUpdate performs UPDATEs using github.com/Masterminds/squirrel.
func SquirrelPreparedUpdate(addresses []*types.Address, tx *sql.Tx) func(b *testing.B) {
	fn := func(b *testing.B) {
		var dbcache *sq.StmtCache
		var err error
		//
		dbcache = sq.NewStmtCache(tx)
		defer dbcache.Clear()
		//
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreUpdate(b)
				//
				query := sq.Update(types.AddressTableName).
					Set("street", address.Street).
					Set("city", address.City).
					Set("state", address.State).
					Set("zip", address.Zip).
					Where(sq.Eq{"pk": address.Id}).
					Suffix("RETURNING modified_tmz").
					RunWith(dbcache).
					PlaceholderFormat(sq.Dollar)
				if err = query.QueryRow().Scan(&address.ModifiedTime); err != nil {
					b.Fatalf("squirrel failed with %v", err.Error())
				}
				//
				address.PostUpdate(b)
			}
		}
	}
	return fn
}
