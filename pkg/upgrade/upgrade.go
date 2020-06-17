package upgrade

import (
	"io"

	"github.com/inconshreveable/go-update"
)

type Upgrade interface {
	Apply(reader io.Reader, opts update.Options) error
}

type DefaultUpgrade struct{}

func (u DefaultUpgrade) Apply(reader io.Reader, opts update.Options) error {
	return update.Apply(reader, opts)
}
