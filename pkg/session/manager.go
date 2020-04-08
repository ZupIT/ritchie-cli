package session

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/crypto/cryptoutil"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

var (
	ErrNoSession = errors.New("please, you need to start a session")
)

type DefaultManager struct {
	sessionFile    string
	passphraseFile string
}

func NewManager(homePath string) DefaultManager {
	return DefaultManager{
		sessionFile:    fmt.Sprintf(sessionFilePattern, homePath),
		passphraseFile: fmt.Sprintf(passphraseFilePattern, homePath),
	}
}

func (d DefaultManager) Create(session Session) error {
	sh := cryptoutil.SumHash(session.Secret)
	passphrase := cryptoutil.EncodeHash(sh)
	session.Secret = passphrase

	if err := fileutil.WriteFile(d.passphraseFile, []byte(passphrase)); err != nil {
		return err
	}

	sb, err := json.Marshal(session)
	if err != nil {
		return err
	}

	hash, err := cryptoutil.SumHashMachine(passphrase)
	if err != nil {
		return err
	}
	cipher := cryptoutil.Encrypt(hash, string(sb))
	if err := fileutil.WriteFilePerm(d.sessionFile, []byte(cipher), 0600); err != nil {
		return err
	}

	return nil
}

func (d DefaultManager) Destroy() error {
	if err := fileutil.RemoveFile(d.sessionFile); err != nil {
		return err
	}

	if err := fileutil.RemoveFile(d.passphraseFile); err != nil {
		return err
	}

	return nil
}

func (d DefaultManager) Current() (Session, error) {
	if !fileutil.Exists(d.sessionFile) || !fileutil.Exists(d.passphraseFile) {
		return Session{}, ErrNoSession
	}

	pb, err := fileutil.ReadFile(d.passphraseFile)
	if err != nil {
		return Session{}, err
	}
	passphrase := string(pb)

	hash, err := cryptoutil.SumHashMachine(passphrase)
	if err != nil {
		return Session{}, err
	}

	sb, err := fileutil.ReadFile(d.sessionFile)
	if err != nil {
		return Session{}, err
	}

	plain := cryptoutil.Decrypt(hash, string(sb))
	session := &Session{}
	if err := json.Unmarshal([]byte(plain), &session); err != nil {
		return Session{}, err
	}

	return *session, nil
}
