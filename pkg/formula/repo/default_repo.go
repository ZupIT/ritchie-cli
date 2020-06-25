package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"

	"github.com/ZupIT/ritchie-cli/pkg/api"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	"github.com/ZupIT/ritchie-cli/pkg/formula"
	"github.com/ZupIT/ritchie-cli/pkg/http/headers"
	"github.com/ZupIT/ritchie-cli/pkg/server"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const (
	// Files
	repositoryConfFilePattern    = "%s/repo/repositories.json"
	repositoryCacheFolderPattern = "%s/repo/cache"
	treeCacheFilePattern         = "%s/repo/cache/%s-tree.json"
	providerPath                 = "%s/repositories"
)

var (
	// Errors
	ErrNoRepoToShow = errors.New("no repositories to show")
)

type RepoManager struct {
	repoFile       string
	cacheFile      string
	homePath       string
	httpClient     *http.Client
	sessionManager session.Manager
	serverFinder   server.Finder
	edition        api.Edition
}

// ByPriority implements sort.Interface for []Repository based on
// the Priority field.
type ByPriority []formula.Repository

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].Priority < a[j].Priority }

func NewSingleRepoManager(homePath string, hc *http.Client, sm session.Manager) RepoManager {
	return RepoManager{
		repoFile:       fmt.Sprintf(repositoryConfFilePattern, homePath),
		cacheFile:      fmt.Sprintf(repositoryCacheFolderPattern, homePath),
		homePath:       homePath,
		httpClient:     hc,
		sessionManager: sm,
		edition:        api.Single,
	}
}

func NewTeamRepoManager(homePath string, serverFinder server.Finder, hc *http.Client, sm session.Manager) RepoManager {
	return RepoManager{

		repoFile:       fmt.Sprintf(repositoryConfFilePattern, homePath),
		cacheFile:      fmt.Sprintf(repositoryCacheFolderPattern, homePath),
		homePath:       homePath,
		serverFinder:   serverFinder,
		httpClient:     hc,
		sessionManager: sm,
		edition:        api.Team,
	}
}

func (dm RepoManager) Add(r formula.Repository) error {
	err := os.MkdirAll(filepath.Dir(dm.cacheFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	lockFile := strings.Replace(dm.repoFile, filepath.Ext(dm.repoFile), ".lock", 1)
	lock := flock.New(lockFile)
	lockCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	locked, err := lock.TryLockContext(lockCtx, time.Second)
	if locked {
		defer func() {
			_ = lock.Unlock()
		}()
	}
	if err != nil {
		return err
	}

	if !fileutil.Exists(dm.repoFile) {
		wb, err := json.Marshal(formula.RepositoryFile{})
		if err != nil {
			return err
		}
		err = fileutil.WriteFile(dm.repoFile, wb)
		if err != nil {
			return err
		}
	}

	rb, err := fileutil.ReadFile(dm.repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var repoFile formula.RepositoryFile
	if err := json.Unmarshal(rb, &repoFile); err != nil {
		return err
	}

	if err := dm.loadTreeFile(r); err != nil {
		return fmt.Errorf("looks like %q is not a valid formula repository or cannot be reached\n", r.TreePath)
	}

	added := false
	for i, v := range repoFile.Values {
		if v.Name == r.Name {
			repoFile.Values[i] = r
			added = true
			break
		}
	}
	if !added {
		repoFile.Values = append(repoFile.Values, r)
	}

	if err := writeFile(repoFile, dm.repoFile, 0644); err != nil {
		return err
	}

	return nil
}

func (dm RepoManager) Update() error {
	f, err := dm.loadReposFromDisk()
	if fileutil.IsNotExistErr(err) || len(f.Values) == 0 {
		return ErrNoRepoToShow
	}

	fmt.Println("Wait while we update your repositories...")
	var wg sync.WaitGroup
	for _, v := range f.Values {
		wg.Add(1)
		go func(v formula.Repository) {
			defer wg.Done()
			if err := dm.loadTreeFile(v); err != nil {
				fmt.Printf("...Unable to get an update from the %q formula repository (%s):\n\t%s\n", v.Name, v.TreePath, err)
			} else {
				fmt.Printf("...Successfully got an update from the %q formula repository\n", v.Name)
			}
		}(v)
	}
	wg.Wait()
	fmt.Println("Done.")

	return nil
}

func (dm RepoManager) Clean(n string) error {
	treeName := fmt.Sprintf(treeCacheFilePattern, dm.homePath, n)

	if err := removeRepoCache(treeName); err != nil {
		return err
	}

	return nil
}

func (dm RepoManager) Delete(name string) error {
	f, err := dm.loadReposFromDisk()
	if fileutil.IsNotExistErr(err) || len(f.Values) == 0 {
		return ErrNoRepoToShow
	}

	l := len(f.Values)
	for i, v := range f.Values {
		if v.Name == name {
			f.Values = append(f.Values[:i], f.Values[i+1:]...)
			break
		}
	}
	if len(f.Values) == l {
		return fmt.Errorf("repository %q not found\n", name)
	}

	if err := writeFile(f, dm.repoFile, 0644); err != nil {
		return err
	}

	treeCacheFile := fmt.Sprintf(treeCacheFilePattern, dm.homePath, name)
	if err := removeRepoCache(treeCacheFile); err != nil {
		return err
	}

	return nil
}

func (dm RepoManager) List() ([]formula.Repository, error) {
	f, err := dm.loadReposFromDisk()

	if fileutil.IsNotExistErr(err) {
		return nil, ErrNoRepoToShow
	}
	if len(f.Values) == 0 {
		return []formula.Repository{}, nil
	}

	sort.Sort(ByPriority(f.Values))

	return f.Values, nil
}

func (dm RepoManager) loadTreeFile(r formula.Repository) error {

	session, err := dm.sessionManager.Current()
	if err != nil {
		return errors.New("error restore current session")
	}
	req, err := http.NewRequest(http.MethodGet, r.TreePath, nil)
	if err != nil {
		return err
	}

	if dm.edition == api.Team {
		req.Header.Set(headers.XOrg, session.Organization)
		req.Header.Set(headers.XRepoName, r.Name)
		req.Header.Set(headers.Authorization, session.AccessToken)
	}
	resp, err := dm.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%d - failed to get index for %s\n", resp.StatusCode, r.TreePath)
	}

	treeFile, err := fileutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	treeCacheFile := fmt.Sprintf(treeCacheFilePattern, dm.homePath, r.Name)
	treeDir := filepath.Dir(treeCacheFile)
	err = fileutil.CreateDirIfNotExists(treeDir, 0755)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(treeCacheFile, treeFile)
	if err != nil {
		return err
	}

	return nil
}

func (dm RepoManager) loadReposFromDisk() (formula.RepositoryFile, error) {
	path := fmt.Sprintf(repositoryConfFilePattern, dm.homePath)
	rf := formula.RepositoryFile{}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return rf, err
	}

	err = json.Unmarshal(b, &rf)
	return rf, err
}

func removeRepoCache(root string) error {
	if _, err := os.Stat(root); fileutil.IsNotExistErr(err) {
		return nil
	} else if err != nil {
		return err
	}
	return os.Remove(root)
}

func writeFile(rf formula.RepositoryFile, path string, perm os.FileMode) error {
	b, err := json.Marshal(rf)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, perm)
}
