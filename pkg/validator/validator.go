package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	urlPatternVersion = "%s/cli-version"
)

//IsValidName validates a name of something
func IsValidName(args []string) error {
	n := len(args)
	if n < 1 {
		return errors.New("requires a name arg")
	}
	if n > 1 {
		return fmt.Errorf("accepts at most 1 arg(s), received %d", n)
	}
	name := args[0]
	if len(name) < 3 {
		return errors.New("name must be at least 3 chars")
	}
	return nil
}

//IsValidLocation validates if location exists
func IsValidLocation(location string) error {
	if !fileutil.Exists(location) {
		return fmt.Errorf("%s is not a valid location", location)
	}
	return nil
}

//HasMinValue validates min value for string
func HasMinValue(str string, min int) error {
	n := len(str)
	if n < min {
		return errors.New("value must contain at least 3 characters")
	}
	return nil
}

//IsValidURL validates the url format
func IsValidURL(value string) error {
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return fmt.Errorf("%s is not a valid url", value)
	}
	return nil
}

//IsValidEmail validate the email format
func IsValidEmail(email string) error {
	rgx := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !rgx.MatchString(email) {
		return fmt.Errorf("%s is not a valid email", email)
	}
	return nil
}

//IsValidVersion Validate version with server
func IsValidVersion(version, org string) {
	url := fmt.Sprintf(urlPatternVersion, env.ServerURL)
	client := &http.Client{Timeout: 2 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("x-org", org)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	cv := struct {
		Url      string `json:"url"`
		Provider string `json:"provider"`
		Version  string `json:"cliversion"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&cv)
	if err != nil {
		return
	}
	if version != cv.Version && len(cv.Version) != 0 {
		log.Printf("[WARNING] Please, update your rit(%s) version to the new release(%s)\n", version, cv.Version)
	}
}
