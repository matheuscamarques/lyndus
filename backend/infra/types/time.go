package types

import (
	"bitbucket.org/lyndus/backend/infra/utils"
	"database/sql/driver"
	"strings"
	"time"
)

type TimeHHMM struct {
	time.Time
}

// UnmarshalJSON Parses the json string in the custom format
func (jt *TimeHHMM) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), `"`)
	nt, err := time.Parse(utils.LAYOUTHOUR, s)
	*jt = TimeHHMM{nt}
	return
}

// MarshalJSON writes a quoted string in the custom format
func (jt TimeHHMM) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.RFC3339Nano)+2)
	b = append(b, '"')
	b = jt.AppendFormat(b, utils.LAYOUTHOUR)
	b = append(b, '"')
	return b, nil
}

func (jt *TimeHHMM) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		jt.Time = t
	}
	return nil
}

// String returns the time in the custom format
func (jt *TimeHHMM) String() string {
	return jt.Format(utils.LAYOUTHOUR) // fmt.Sprintf("%q", jt.Format(utils.LAYOUTHOUR))

}

func (jt TimeHHMM) Value() (driver.Value, error) {
	return jt.String(), nil
}

func (jt *TimeHHMM) ParseISO(s string) (err error) {
	nt, err := time.Parse(utils.LAYOUTHOUR, s)
	nt = time.Date(0, 1, 1, nt.Hour(), nt.Minute(), 0, 0, &time.Location{})
	*jt = TimeHHMM{nt}
	return err
}
