package stdin

import (
	"encoding/json"
	"os"
)

// ReadJson reads the json from stdin inputs
func ReadJson(v interface{}) error {
	return json.NewDecoder(os.Stdin).Decode(v)
}

