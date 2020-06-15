package version

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type StubResolverWithSameVersions struct {
}

func (r StubResolverWithSameVersions) GetCurrentVersion() (string, error) {
	return "1.0.0", nil
}

func (r StubResolverWithSameVersions) GetStableVersion() (string, error) {
	return "1.0.0", nil
}

type StubResolverWithDifferentVersions struct {
}

func (r StubResolverWithDifferentVersions) GetCurrentVersion() (string, error) {
	return "1.0.0", nil
}

func (r StubResolverWithDifferentVersions) GetStableVersion() (string, error) {
	return "1.0.1", nil
}

type StubResolverWithErrorOnGetCurrentVersion struct {
}

func (r StubResolverWithErrorOnGetCurrentVersion) GetCurrentVersion() (string, error) {
	return "", errors.New("some error")
}

func (r StubResolverWithErrorOnGetCurrentVersion) GetStableVersion() (string, error) {
	return "1.0.1", nil
}

type StubResolverWithErrorOnGetStableVersion struct {
}

func (r StubResolverWithErrorOnGetStableVersion) GetCurrentVersion() (string, error) {
	return "1.0.0", nil
}

func (r StubResolverWithErrorOnGetStableVersion) GetStableVersion() (string, error) {
	return "1.0.1", errors.New("some error")
}

func TestVerifyNewVersion(t *testing.T) {

	testCases := []struct {
		name           string
		resolver       Resolver
		expectedResult string
	}{
		{
			name:           "Should not print warning",
			resolver:       StubResolverWithSameVersions{},
			expectedResult: "",
		},
		{
			name:           "Should print warning",
			resolver:       StubResolverWithDifferentVersions{},
			expectedResult: fmt.Sprintf(prompt.Warning, MsgRitUpgrade),
		},
		{
			name:           "Should not print on error in GetCurrentVersion",
			resolver:       StubResolverWithErrorOnGetCurrentVersion{},
			expectedResult: "",
		},
		{
			name:           "Should not print on error in GetStableVersion",
			resolver:       StubResolverWithErrorOnGetStableVersion{},
			expectedResult: "",
		},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {

			buffer := &bytes.Buffer{}

			VerifyNewVersion(tCase.resolver, buffer)

			result := buffer.String()

			assertEquals(result, tCase.expectedResult, t)
		})
	}
}

func TestGetStableVersion(t *testing.T) {

	t.Run("Should get stableVersion", func(t *testing.T) {
		expectedResult := "1.0.0"

		mockHttp := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expectedResult + "\n"))
		}))

		result, err := DefaultVersionResolver{StableVersionUrl: mockHttp.URL}.GetStableVersion()
		if err != nil {
			t.Errorf("fail Err:%s\n", err)
		}
		assertEquals(expectedResult, result, t)
	})

	t.Run("Should return err when http.get fail", func(t *testing.T) {
		_, err := DefaultVersionResolver{}.GetStableVersion()
		if err == nil {
			t.Fatalf("Should return err.")
		}
	})

}

func TestGetCurrentVersion(t *testing.T) {
	t.Run("Should Return the Current Version", func(t *testing.T) {
		currentVersion := "0.0.1"
		resolver := DefaultVersionResolver{CurrentVersion: currentVersion}
		result, _ := resolver.GetCurrentVersion()
		assertEquals(currentVersion, result, t)
	})
}

func assertEquals(expected string, result string, t *testing.T) {
	if expected != result {
		t.Helper()
		t.Errorf("\nExpected: %s\nbut was:%s\n", expected, result)
	}
}
