package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/inconshreveable/go-update"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

const (
	stableVersion  = "1.0.0"
	currentVersion = "1.1.0"
)

var (
	stubUpgradeUtilApplyExecutions     = 0
	stubUpgradeUtilFailApplyExecutions = 0
)

type stubResolver struct{}

func (r stubResolver) GetCurrentVersion() (string, error) {
	return currentVersion, nil
}

func (r stubResolver) GetStableVersion() (string, error) {
	return stableVersion, nil
}

type stubResolverWithError struct{}

func (r stubResolverWithError) GetCurrentVersion() (string, error) {
	return currentVersion, nil
}

func (r stubResolverWithError) GetStableVersion() (string, error) {
	return "", errors.New("some Error")
}

type StubUpgradeUtil struct{}

func (u StubUpgradeUtil) Apply(reader io.Reader, opts update.Options) error {
	stubUpgradeUtilApplyExecutions++
	return nil
}

type StubUpgradeUtilFail struct{}

func (u StubUpgradeUtilFail) Apply(reader io.Reader, opts update.Options) error {
	stubUpgradeUtilFailApplyExecutions++
	return errors.New("Some Error")
}

func TestGetUpgradeUrlSingle(t *testing.T) {
	result := GetUpgradeUrl(api.Single, stubResolver{})
	expected := fmt.Sprintf(upgradeUrlFormat, stableVersion, runtime.GOOS, api.Single)
	if runtime.GOOS == "windows" {
		expected += ".exe"
	}
	assertEquals(expected, result, t)
}

func TestGetUpgradeUrlTeam(t *testing.T) {
	result := GetUpgradeUrl(api.Team, stubResolver{})
	expected := fmt.Sprintf(upgradeUrlFormat, stableVersion, runtime.GOOS, api.Team)
	if runtime.GOOS == "windows" {
		expected += ".exe"
	}
	assertEquals(expected, result, t)
}

func TestGetUpgradeWithError(t *testing.T) {
	result := GetUpgradeUrl(api.Team, stubResolverWithError{})
	expected := ""
	assertEquals(expected, result, t)
}

func TestNewUpgradeCmd(t *testing.T) {

	stubUpgradeUtilApplyExecutions = 0

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer mockServer.Close()

	upgradeCmd := NewUpgradeCmd(mockServer.URL, StubUpgradeUtil{})

	if upgradeCmd == nil {
		t.Errorf("Expeceted to build a UpgradeCmd")
	}

	err := upgradeCmd.Execute()

	if stubUpgradeUtilApplyExecutions != 1 {
		t.Errorf("Expected 1 executions of StubUpgradeUtil.apply")
	}

	if err != nil {
		t.Errorf("Error to execute upgradeCmd")
	}

}

func TestNewUpgradeCmdWrongUrl(t *testing.T) {

	upgradeUrl := ""
	stubUpgradeUtilApplyExecutions = 0

	upgradeCmd := NewUpgradeCmd(upgradeUrl, StubUpgradeUtil{})

	if upgradeCmd == nil {
		t.Errorf("Expeceted to build a UpgradeCmd")
	}

	err := upgradeCmd.Execute()

	if err == nil {
		t.Errorf("Expected Error")
	}

	if stubUpgradeUtilApplyExecutions != 0 {
		t.Errorf("Shloud not call upgradeUtil.")
	}

}

func TestNewUpgradeCmdFailToGet(t *testing.T) {

	upgradeUrl := "someUrl"
	stubUpgradeUtilApplyExecutions = 0

	upgradeCmd := NewUpgradeCmd(upgradeUrl, StubUpgradeUtil{})

	if upgradeCmd == nil {
		t.Errorf("Expeceted to build a UpgradeCmd")
	}

	err := upgradeCmd.Execute()

	if err == nil {
		t.Errorf("Expected Error")
	}

	if stubUpgradeUtilApplyExecutions != 0 {
		t.Errorf("Shloud not call upgradeUtil.")
	}

}

func TestNewUpgradeCmdNotFoundToGet(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer mockServer.Close()
	stubUpgradeUtilApplyExecutions = 0

	upgradeCmd := NewUpgradeCmd(mockServer.URL, StubUpgradeUtil{})

	if upgradeCmd == nil {
		t.Errorf("Expeceted to build a UpgradeCmd")
	}

	err := upgradeCmd.Execute()

	if err == nil {
		t.Errorf("Expected Error")
	}

	if stubUpgradeUtilApplyExecutions != 0 {
		t.Errorf("Shloud not call upgradeUtil.")
	}

}

func TestNewUpgradeCmdFailToApplyUpgrade(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer mockServer.Close()
	stubUpgradeUtilFailApplyExecutions = 0

	upgradeCmd := NewUpgradeCmd(mockServer.URL, StubUpgradeUtilFail{})

	if upgradeCmd == nil {
		t.Errorf("Expeceted to build a UpgradeCmd")
	}

	err := upgradeCmd.Execute()

	if err == nil {
		t.Errorf("Expected Error")
	}

	if stubUpgradeUtilFailApplyExecutions != 1 {
		t.Errorf("Expected 1 executions of StubUpgradeUtil.apply")
	}

}

func assertEquals(expected string, result string, t *testing.T) {
	if expected != result {
		t.Helper()
		t.Errorf("\nExpected: %s\nbut was:%s\n", expected, result)
	}
}
