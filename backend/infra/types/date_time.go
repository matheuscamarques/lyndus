package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/lyndus/backend/infra/utils"
)

type DateTime struct {
	time.Time
}

// UnmarshalJSON Parses the json string in the custom format
func (dt *DateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(utils.LAYOUTDATETIME1, s)
	if err != nil {
		nt, err = time.Parse(utils.LAYOUTDATETIME2, s)
	}
	*dt = DateTime{nt}
	return err
}

// MarshalJSON writes a quoted string in the custom format
func (dt DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.RFC3339Nano)+2)
	b = append(b, '"')
	b = dt.AppendFormat(b, utils.LAYOUTDATETIME1)
	b = append(b, '"')
	return b, nil

}

func (dt *DateTime) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		dt.Time = t
	}
	return nil
}

func (dt *DateTime) ParseISO(s string) (err error) {
	nt, err := time.Parse(utils.LAYOUTDATETIME1, s)
	if err != nil {
		nt, err = time.Parse(utils.LAYOUTDATETIME2, s)
	}
	*dt = DateTime{nt}
	return err
}

// String returns the time in the custom format
func (dt *DateTime) String() string {
	return fmt.Sprintf("%q", dt.Format(utils.LAYOUTDATETIME1))
}

func (dt *DateTime) SetNow(locSTR string) (err error) {
	location, err := time.LoadLocation(locSTR)
	if err != nil {
		return
	}

	dt.Time = time.Now().In(location)
	dt.Time = time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour(), dt.Minute(), 0, 0, &time.Location{})
	return
}

func (dt DateTime) Value() (driver.Value, error) {
	return dt.String(), nil
}

type DateTimeTZ struct{ time.Time }

// UnmarshalJSON Parses the json string in the custom format
func (dt *DateTimeTZ) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(time.RFC3339, s)
	*dt = DateTimeTZ{nt}
	return
}

// MarshalJSON writes a quoted string in the custom format
func (dt DateTimeTZ) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.RFC3339Nano)+2)
	b = append(b, '"')
	b = dt.AppendFormat(b, time.RFC3339)
	b = append(b, '"')
	return b, nil
}

// String returns the time in the custom format
func (dt *DateTimeTZ) String() string {
	return dt.Format(time.RFC3339) //fmt.Sprintf("%q", t.Format(time.RFC3339))
}

func (dt *DateTimeTZ) ParseISO(date string, hour string) (err error) {
	//date is 2021-08-19
	//hour is 19:32
	t, err := time.Parse("2006-01-02T15:04", date+"T"+hour)
	if err != nil {
		return
	}
	*dt = DateTimeTZ{t}
	return
}
