package session

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ZupIT/ritchie-cli/pkg/crypto/cryptoutil"
	"github.com/ZupIT/ritchie-cli/pkg/stream"
)

var (
	ErrNoSession = errors.New("please, you need to start a session")
)

type DefaultManager struct {
	sessionFile    string
	passphraseFile string
	file           stream.FileWriteReadExistRemover
}

func NewManager(homePath string, file stream.FileWriteReadExistRemover) DefaultManager {
	return DefaultManager{
		sessionFile:    fmt.Sprintf(sessionFilePattern, homePath),
		passphraseFile: fmt.Sprintf(passphraseFilePattern, homePath),
		file:           file,
	}
}

func (d DefaultManager) Create(session Session) error {
	sh := cryptoutil.SumHash(session.Secret)
	passphrase := cryptoutil.EncodeHash(sh)
	session.Secret = passphrase

	if err := d.file.Write(d.passphraseFile, []byte(passphrase)); err != nil {
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
	if err := d.file.Write(d.sessionFile, []byte(cipher)); err != nil {
		return err
	}

	return nil
}

func (d DefaultManager) Destroy() error {
	if err := d.file.Remove(d.sessionFile); err != nil {
		return err
	}

	if err := d.file.Remove(d.passphraseFile); err != nil {
		return err
	}

	return nil
}

func (d DefaultManager) Current() (Session, error) {
	if !d.file.Exists(d.sessionFile) || !d.file.Exists(d.passphraseFile) {
		return Session{}, ErrNoSession
	}

	pb, err := d.file.Read(d.passphraseFile)
	if err != nil {
		return Session{}, err
	}
	passphrase := string(pb)

	hash, err := cryptoutil.SumHashMachine(passphrase)
	if err != nil {
		return Session{}, err
	}

	sb, err := d.file.Read(d.sessionFile)
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
