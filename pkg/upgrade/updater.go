package upgrade

import (
	"io"

	"github.com/inconshreveable/go-update"
)

type Updater interface {
	Apply(reader io.Reader, opts update.Options) error
}

type DefaultUpdater struct{}

func (u DefaultUpdater) Apply(reader io.Reader, opts update.Options) error {
	return update.Apply(reader, opts)
}
