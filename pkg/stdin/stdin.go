package stdin

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

const (
	MsgInvalidInput = "The STDIN inputs weren't informed correctly. Check the JSON used to execute the command."
)

// ReadJson reads the json from stdin inputs
func ReadJson(reader io.Reader, v interface{}) error {
	return json.NewDecoder(reader).Decode(v)
}

// This function aims to facilitate stdin manipulations especially on tests
// inputs:
// 		content string: the string content for the stdin
// outputs:
//		tmpfile File: the tmp file read by stdin
// 		oldstdin File: the original stdin being overwritten
// NOTICE: do not forget to remove the tmpfile and restore the original stdin, use the commands
// 		defer func() { os.Stdin = oldStdin }()
// 		defer os.Remove(tmpfile.Name())
func WriteToStdin(content string) (tmpfile *os.File, oldStdin *os.File, error error) {
	fileContent := []byte(content)
	tmpfile, err := ioutil.TempFile("", "stdin")
	if err != nil {
		return nil, nil, err
	}

	if _, err := tmpfile.Write(fileContent); err != nil {
		return nil, nil, err
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		return nil, nil, err
	}

	oldstdin := os.Stdin
	os.Stdin = tmpfile

	return tmpfile, oldstdin, nil
}
