package types

import (
	"encoding/json"
	"io"
)

type Object struct{}

func (o Object) AssignBody(target interface{}, source io.Reader) error {
	return json.NewDecoder(source).Decode(target)
}
