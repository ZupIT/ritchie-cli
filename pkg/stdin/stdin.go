package stdin

import (
	"encoding/json"
	"bufio"
	"os"
	"strings"
)

// Args type that represents parsed args from stdin
type Args map[string]string

// ReadJson reads the json from stdin inputs
func ReadJson(v interface{}) error {
	return json.NewDecoder(os.Stdin).Decode(v)
}

// Parse reads the input from stdin and parse the args to map[name]value
func Parse() (Args, error) {
	reader := bufio.NewReader(os.Stdin)
	data, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	data = strings.TrimSpace(data)
	ss := strings.Split(data, "@in:")
	args := make(Args, len(ss))
	for _, s := range ss {
		if s == "" || s == " " {
			continue
		}

		kv := strings.Split(s, "=")
		k := kv[0]
		v := kv[1]
		args[k] = v
	}

	return args, nil
}

// TODO : formulas STDIN

// r := make(map[string]interface{})

// err := stdin.ReadJson(&r)
// if err != nil {
//	fmt.Println("The stdin inputs weren't informed correctly. Check the JSON used to execute the command.")
//	return err
//}

// fmt.Println("Map:", r)

// TODO : check config.json inputs with map

// if err := execute.FormulaCommand(nil); err != nil {
//	return err
//}