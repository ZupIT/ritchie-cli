package stdin

import (
	"encoding/json"
	"errors"
	"io"
)

var ErrInvalidInput = errors.New("the STDIN inputs weren't informed correctly. Check the JSON used to execute the command")

// ReadJson reads the json from stdin inputs
func ReadJson(reader io.Reader, v interface{}) error {
	if err := json.NewDecoder(reader).Decode(v); err != nil {
		return ErrInvalidInput
	}

	return nil
}
