package types

import (
	"encoding/json"

	"github.com/nofeaturesonlybugs/sqlhbenchmarks/data"
)

// SaleRecords is initialized when importing the package.
var SaleRecords []*SaleReport
var SizeSaleRecords int

// AddressRecords is initialized when importing the package.
var AddressRecords []*Address
var SizeAddressRecords int

func init() {
	{
		dest := []*SaleReport{}
		err := json.Unmarshal([]byte(data.JsonSales), &dest)
		if err != nil {
			panic("init sales json with " + err.Error())
		}
		SaleRecords = dest
		SizeSaleRecords = len(SaleRecords)
	}
	//
	{
		dest := []*Address{}
		err := json.Unmarshal([]byte(data.JsonAddresses), &dest)
		if err != nil {
			panic("init addresses json with " + err.Error())
		}
		AddressRecords = dest
		SizeAddressRecords = len(AddressRecords)
	}
}
