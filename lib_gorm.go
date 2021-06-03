package sqlhbenchmarks

import (
	"testing"

	"github.com/nofeaturesonlybugs/sqlhbenchmarks/types"

	"gorm.io/gorm"
)

// GORMSelect selects records using GORM.
func GORMSelect(limit int, db *gorm.DB) func(*testing.B) {
	fn := func(b *testing.B) {
		var dest []*types.Address
		var result *gorm.DB
		//
		for k := 0; k < b.N; k++ {
			result = db.Limit(limit).Find(&dest)
			if result.Error != nil {
				b.Fatalf("gorm failed with %v", result.Error.Error())
			}
		}
	}
	return fn
}

// GORMInsert performs INSERTs using GORM.
func GORMInsert(addresses []*types.Address, db *gorm.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var result *gorm.DB
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreInsert(b)
				//
				result = db.Create(address)
				if result.Error != nil {
					b.Fatalf("gorm failed with %v", result.Error.Error())
				}
				// address.PostInsert(b) // TODO CreatedAt, ModifiedAt not working with our "stacked" model type.
			}
		}
	}
	return fn
}

// GORMPreparedInsert performs INSERTs using GORM.
func GORMPreparedInsert(addresses []*types.Address, db *gorm.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var result *gorm.DB
		b.StopTimer()
		copies := make([]*types.Address, len(addresses))
		b.StartTimer()
		for k := 0; k < b.N; k++ {
			b.StopTimer()
			for k, v := range addresses {
				copies[k] = &types.Address{
					Street: v.Street,
					City:   v.City,
					State:  v.State,
					Zip:    v.Zip,
				}
			}
			b.StartTimer()
			result = db.Create(copies)
			if result.Error != nil {
				b.Fatalf("gorm failed with %v", result.Error.Error())
			}
		}
	}
	return fn
}

// GORMUpdate performs UPDATEs using GORM.
func GORMUpdate(addresses []*types.Address, db *gorm.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var result *gorm.DB
		for k := 0; k < b.N; k++ {
			for _, address := range addresses {
				address.PreUpdate(b)
				//
				result = db.Save(address)
				if result.Error != nil {
					b.Fatalf("gorm failed with %v", result.Error.Error())
				}
				// address.PostUpdate(b) // TODO CreatedAt, ModifiedAt not working with our "stacked" model type.
			}
		}
	}
	return fn
}

// GORMPreparedUpdate performs INSERTs using GORM.
func GORMPreparedUpdate(addresses []*types.Address, db *gorm.DB) func(b *testing.B) {
	fn := func(b *testing.B) {
		var result *gorm.DB
		for k := 0; k < b.N; k++ {
			result = db.Save(addresses)
			if result.Error != nil {
				b.Fatalf("gorm failed with %v", result.Error.Error())
			}
		}
	}
	return fn
}
