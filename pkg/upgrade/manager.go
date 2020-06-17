package upgrade

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/inconshreveable/go-update"
)

type Manager interface {
	Run(upgradeUrl string) error
}

type DefaultManager struct {
	Updater
}

func (m DefaultManager) Run(upgradeUrl string) error {
	if upgradeUrl == "" {
		return errors.New("fail to resolve upgrade url")
	}

	resp, err := http.Get(upgradeUrl)
	if err != nil {
		return errors.New("fail to download stable version")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("fail to download stable version status:%d", resp.StatusCode)
	}

	err = m.Updater.Apply(resp.Body, update.Options{})
	if err != nil {
		return fmt.Errorf("fail to upgrade")
	}
	return nil
}
