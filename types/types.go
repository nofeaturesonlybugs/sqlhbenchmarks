package types

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nofeaturesonlybugs/set"
	"github.com/nofeaturesonlybugs/sqlh/grammar"
	"github.com/nofeaturesonlybugs/sqlh/model"
)

// Model table names.
var AddressTableName = "sqlh_addresses"

// MockRows is the interface for mocking rows.
type MockRows interface {
	MockRows(int) *sqlmock.Rows
}

// TestModel represents models that we test.  By implementing this interface our test code
// can be more uniform.
type TestModel interface {
	PreInsert(b *testing.B)
	PostInsert(b *testing.B)
	PreUpdate(b *testing.B)
	PostUpdate(b *testing.B)
}

// NewMapper returns an appropriate *set.Mapper for the types in this package.
func NewMapper() *set.Mapper {
	rv := &set.Mapper{
		TreatAsScalar: set.NewTypeList(Time{}),
		Join:          "_",
		Tags:          []string{"db", "json"},
	}
	return rv
}

// NewModels returns a model.Models for the types in this package.
func NewModels(grammar grammar.Grammar) *model.Models {
	rv := &model.Models{
		Mapper:  NewMapper(),
		Grammar: grammar,
	}
	rv.Register(&Address{}, model.TableName(AddressTableName))
	return rv
}

// Address model to represent an address.
//
// Methods like Pop, Push, PreInsert, TestInsert, and TestUpdate are present to simplify test code and are not
// required by the model package.
type Address struct {
	Id           int    `json:"id" db:"pk" model:"key,auto" gorm:"column:pk;primaryKey"`
	CreatedTime  Time   `json:"created_time" db:"created_tmz" model:"inserted" gorm:"-"`
	ModifiedTime Time   `json:"modified_time" db:"modified_tmz" model:"inserted,updated" gorm:"-"`
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	//
	pushModified Time
}

func (me *Address) MockRows(n int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"pk", "created_tmz", "modified_tmz",
		"street", "city", "state", "zip",
	})
	for k := 0; k < n; k++ {
		j := AddressRecords[k%SizeAddressRecords]
		rows.AddRow(
			j.Id, j.CreatedTime, j.ModifiedTime,
			j.Street, j.City, j.State, j.Zip,
		)
	}
	return rows
}

func (me *Address) PreInsert(b *testing.B) {
	me.Id, me.CreatedTime, me.ModifiedTime = 0, ZeroTime, ZeroTime
}

func (me *Address) PostInsert(b *testing.B) {
	if me.Id <= 0 {
		b.Fatalf("%T.Id not updated", me)
	} else if me.CreatedTime.IsZero() {
		b.Fatalf("%T.CreatedTime not updated", me)
	} else if me.ModifiedTime.IsZero() {
		b.Fatalf("%T.ModifiedTime not updated", me)
	}
}

func (me *Address) PreUpdate(b *testing.B) {
	me.pushModified = me.ModifiedTime
}

func (me *Address) PostUpdate(b *testing.B) {
	if me.pushModified.Equal(me.ModifiedTime.Time) {
		b.Fatalf("%T address not updated.", me)
	}
	me.ModifiedTime = me.pushModified
}

// TableName overrides the table name used by User to `profiles`.
func (me *Address) TableName() string {
	return AddressTableName
}

// SaleReport is a SELECT destination; it does not represent models.
type SaleReport struct {
	Id                 int    `json:"id" db:"pk"`
	CreatedTime        string `json:"created_time" db:"created_tmz"`
	ModifiedTime       string `json:"modified_time" db:"modified_tmz"`
	Price              int    `json:"price" db:"price"`
	Quantity           int    `json:"quantity" db:"quantity"`
	Total              int    `json:"total" db:"total"`
	CustomerId         int    `json:"customer_id" db:"customer_id"`
	CustomerFirst      string `json:"customer_first" db:"customer_first"`
	CustomerLast       string `json:"customer_last" db:"customer_last"`
	VendorId           int    `json:"vendor_id" db:"vendor_id"`
	VendorName         string `json:"vendor_name" db:"vendor_name"`
	VendorDescription  string `json:"vendor_description" db:"vendor_description"`
	VendorContactId    int    `json:"vendor_contact_id" db:"vendor_contact_id"`
	VendorContactFirst string `json:"vendor_contact_first" db:"vendor_contact_first"`
	VendorContactLast  string `json:"vendor_contact_last" db:"vendor_contact_last"`
}

func (me *SaleReport) MockRows(n int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"pk", "created_tmz", "modified_tmz",
		"price", "quantity", "total",
		"customer_id", "customer_first", "customer_last",
		"vendor_id", "vendor_name", "vendor_description",
		"vendor_contact_id", "vendor_contact_first", "vendor_contact_last",
	})
	for k := 0; k < n; k++ {
		j := SaleRecords[k%SizeSaleRecords]
		rows.AddRow(
			j.Id, j.CreatedTime, j.ModifiedTime,
			j.Price, j.Quantity, j.Total,
			j.CustomerId, j.CustomerFirst, j.CustomerLast,
			j.VendorId, j.VendorName, j.VendorDescription,
			j.VendorContactId, j.VendorContactFirst, j.VendorContactLast,
		)
	}
	return rows
}
