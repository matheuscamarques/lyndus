package types

import (
	"bitbucket.org/lyndus/backend/infra/utils"
	"database/sql/driver"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

// UnmarshalJSON Parses the json string in the custom format
func (jt *Date) UnmarshalJSON(b []byte) (err error) {

	s := strings.Trim(string(b), `"`)

	if len(s) == 0 {
		return err
	}
	jt.Time, err = time.Parse(utils.LAYOUTDATE2, s)
	if err != nil {
		jt.Time, err = time.Parse(utils.LAYOUTDATE, s)
	}
	return err
}

// MarshalJSON writes a quoted string in the custom format
func (jt Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.RFC3339Nano)+2)
	b = append(b, '"')
	b = jt.AppendFormat(b, utils.LAYOUTDATE2)
	b = append(b, '"')
	return b, nil
}

func (jt *Date) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		jt.Time = t
	}
	return nil
}

// String returns the time in the custom format
func (jt *Date) String() string {
	return jt.Format(utils.LAYOUTDATE)
}

func (jt *Date) SetNow() {
	jt.Time = time.Now()
	jt.Time = time.Date(jt.Year(), jt.Month(), jt.Day(), 0, 0, 0, 0, &time.Location{})
}

func (jt Date) Value() (driver.Value, error) {
	return jt.String(), nil
}

func (jt *Date) ParseISO(s string) (err error) {
	nt, err := time.Parse(utils.LAYOUTDATE, s)
	nt = time.Date(nt.Year(), nt.Month(), nt.Day(), 0, 0, 0, 0, &time.Location{})
	*jt = Date{nt}
	return
}
