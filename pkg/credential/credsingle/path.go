package credsingle

import (
	"fmt"
)

const (
	pathPattern     = "%s/credentials/%s"
	credFilePattern = "%s/%s"
)

func Dir(homePath, ctx string) string {
	return fmt.Sprintf(pathPattern, homePath, ctx)
}

func File(homePath, ctx, provider string) string {
	return fmt.Sprintf(credFilePattern, Dir(homePath, ctx), provider)
}
