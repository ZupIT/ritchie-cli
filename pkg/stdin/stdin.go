package stdin

import (
	"bufio"
	"os"
	"strings"
)

// Args type that represents parsed args from stdin
type Args map[string]string

// Parse reads the input from stdin and parse the args to map[name]value
func Parse() (Args, error) {
	reader := bufio.NewReader(os.Stdin)
	data, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	//data = strings.TrimSpace(data)

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
