package stdin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"testing"
)

const msg = "read stdin test"

type TestReader struct {
	Test string `json:"test"`
}

func TestReadJson(t *testing.T) {

	// Convert interface to Json for test
	i := TestReader{Test: msg}
	jsonData, _ := json.Marshal(i)

	// Insert Json inside a new Reader (simulating os.Stdin)
	var stdin bytes.Buffer
	stdin.Write(jsonData)
	reader := bufio.NewReader(&stdin)

	tr := TestReader{}

	// ReadJson through Reader and convert to chosen interface
	err := ReadJson(reader, &tr); if err != nil {
		t.Errorf("Got error %v", err)
	}

	// Assert the decoder result is the initial message
	if msg != tr.Test {
		t.Errorf("Expected : %v but got %v", msg, tr.Test)
	}
}