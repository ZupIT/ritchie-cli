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

package prompt

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
)

// NewError returns new error with red message
func NewError(text string) error {
	return errors.New(Red(text))
}

// Red returns a red string
func Red(text string) string {
	return color.FgRed.Render(text)
}

// Error is a Println with red message
func Error(text string) {
	fmt.Println(Red(text))
}

func Green(text string) string {
	return color.Success.Render(text)
}
func Success(text string) {
	fmt.Println(Green(text))
}

func Bold(text string) string {
	return color.Bold.Render(text)
}
func Info(text string) {
	fmt.Println(Bold(text))
}

func Yellow(text string) string {
	return color.Warn.Render(text)
}
func Warning(text string) {
	color.Warn.Println(text)
}

func Cyan(text string) string {
	return color.Cyan.Render(text)
}
