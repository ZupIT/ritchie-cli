/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package stdin

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

var ErrInvalidInput = errors.New("the STDIN inputs weren't informed correctly. Check the JSON used to execute the command")

// ReadJson reads the json from stdin inputs
func ReadJson(reader io.Reader, v interface{}) error {
	if err := json.NewDecoder(reader).Decode(v); err != nil {
		return ErrInvalidInput
	}

	return nil
}

// func ExistsEntry(reader io.Reader) bool {
// 	type Message struct {
// 		Name, Text string
// 	}
// 	var m Message

// 	if err := ReadJson(reader, &m); err != nil {
// 		return false
// 	}
// 	return true
// }

func reading(r io.Reader, w io.Writer, v interface{}) {
	nr := bufio.NewReader(r)
	b, e := nr.Peek(1)
	fmt.Println("peek v")
	fmt.Println(b, e)
	fmt.Println("peek ^")
	// buf := make([]byte, 4)
	// for {
	// 	n, err := r.Read(buf)
	// 	fmt.Println(n, err, buf[:n])
	// 	if err == io.EOF {
	// 		break
	// 	}
	// }
}

func ExistsEntry(reader io.Reader, cmd *cobra.Command) bool {
	init := []string{""}

	go reading(cmd.InOrStdin(), cmd.OutOrStdout(), &init)

	fmt.Println("init: ", init)

	tempo := 10
	i := 1
	for i < tempo {
		time.Sleep(1 * time.Second)
		fmt.Println(".")

		i++
	}

	return true
}
