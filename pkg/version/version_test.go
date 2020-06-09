package version

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ZupIT/ritchie-cli/pkg/prompt"
)

type StubResolverWithSameVersions struct {
}

func (r StubResolverWithSameVersions) getCurrentVersion() string {
	return "1.0.0"
}

func (r StubResolverWithSameVersions) getStableVersion() string {
	return "1.0.0"
}

type StubResolverWithDifferentVersions struct {
}

func (r StubResolverWithDifferentVersions) getCurrentVersion() string {
	return "1.0.0"
}

func (r StubResolverWithDifferentVersions) getStableVersion() string {
	return "1.0.1"
}


func TestVerifyNewVersion(t *testing.T) {

	testCases := []struct {
		name           string
		resolver       Resolver
		expectedResult string
	}{
		{
			name: "Should not print warning",
			resolver: StubResolverWithSameVersions{},
			expectedResult: "",
		},
		{
			name: "Should print warning",
			resolver: StubResolverWithDifferentVersions{},
			expectedResult: fmt.Sprintf(prompt.Warning,MsgRitUpgrade),
		},

	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {

			buffer := &bytes.Buffer{}

			VerifyNewVersion(tCase.resolver, buffer)

			result := buffer.String()

			if result != tCase.expectedResult {
				t.Errorf("\nExpected: %s\nbut was:%s\n", result, tCase.expectedResult)
			}
		})
	}
}
