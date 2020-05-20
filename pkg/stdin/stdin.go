package stdin

import (
	"encoding/json"
	"io"
)

// ReadJson reads the json from stdin inputs
func ReadJson(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}

