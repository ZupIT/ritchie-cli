package stdin

import (
	"encoding/json"
	"io"
)

var (
	MsgInvalidInput = "The STDIN inputs weren't informed correctly. Check the JSON used to execute the command."
)

// ReadJson reads the json from stdin inputs
func ReadJson(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}
