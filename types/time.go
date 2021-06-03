package types

import (
	"database/sql/driver"
	"time"

	"github.com/nofeaturesonlybugs/errors"
)

var ZeroTime Time

// Time overloads time.Time so we can use the same type in our models for CreatedTime and ModifiedTime
// and use it in multiple database drivers.
type Time struct {
	time.Time
}

// Scan implements the Scanner interface.
func (me *Time) Scan(value interface{}) error {
	var err error
	switch v := value.(type) {
	case time.Time:
		me.Time = v

	case string:
		if me.Time, err = time.ParseInLocation("2006-01-02 15:04:05", v, time.UTC); err != nil {
			err = errors.Go(err)
		}

	case int64:
		me.Time = time.Unix(v, 0)

	default:
		err = errors.Errorf("%T unsupported for SqliteTime", value)
	}
	return err
}

// Value implements the driver Valuer interface.
func (me *Time) Value() (driver.Value, error) {
	return me.Time.Unix(), nil
}

func (*Time) GormDataType() string {
	return "TIME"
}
