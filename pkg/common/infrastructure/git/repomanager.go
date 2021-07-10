package git

import (
	"bytes"
	"math"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"patcher/pkg/common/infrastructure/executor"
)

func NewRepoManager(repoPath string, gitExecutor Executor) RepoManager {
	return &repoManager{
		repoPath: repoPath,
		executor: gitExecutor,
	}
}

type RepoManager interface {
	Checkout(branch string) error
	ForceCheckout(branch string) error
	Fetch() error
	FetchAll() error
	ApplyPatch(patchContent []byte) error
	HardReset() error

	RemoteBranches() ([]string, error)
	// ListChangedFiles returns slice of changed files
	ListChangedFiles() ([]string, error)
	GetCurrentChanges(withCached bool) ([]byte, error)
	RemoteProjectName() (string, error)
}

type repoManager struct {
	repoPath string
	executor Executor
}

func (repo *repoManager) Checkout(branch string) error {
	return repo.run("checkout", branch)
}

func (repo *repoManager) ForceCheckout(branch string) error {
	return repo.run("checkout", "-f", branch)
}

func (repo *repoManager) Fetch() error {
	return repo.run("fetch")
}

func (repo *repoManager) FetchAll() error {
	return repo.run("fetch", "--all")
}

func (repo *repoManager) ApplyPatch(patchContent []byte) error {
	return repo.runWithOpts([]string{"apply"}, executor.WithStdin(bytes.NewBuffer(patchContent)))
}

func (repo *repoManager) HardReset() error {
	return repo.runWithOpts([]string{"reset", "--hard"})
}

//nolint:prealloc
func (repo *repoManager) RemoteBranches() ([]string, error) {
	output, err := repo.output("remote", "-v")
	if err != nil {
		return nil, err
	}

	reg := regexp.MustCompile(`(^.+?)\s`)

	var branches []string

	for i, s := range strings.Split(string(output), "\n") {
		if math.Mod(float64(i), 2) != 0 {
			continue
		}
		if s == "" {
			continue
		}
		branches = append(branches, strings.TrimSpace(reg.FindString(s)))
	}
	return branches, nil
}

func (repo *repoManager) ListChangedFiles() ([]string, error) {
	output, err := repo.output("status", "-s")
	if err != nil {
		return nil, err
	}

	const unversionedPrefix = "??"

	var result []string

	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, unversionedPrefix) {
			result = append(result, line)
		}
	}

	return result, nil
}

func (repo *repoManager) GetCurrentChanges(withCached bool) ([]byte, error) {
	args := []string{"diff"}
	if withCached {
		args = append(args, "--cached")
	}

	return repo.output(args...)
}

func (repo *repoManager) RemoteProjectName() (string, error) {
	remoteNameBytes, err := repo.output("remote")
	if err != nil {
		return "", errors.WithStack(err)
	}

	if len(remoteNameBytes) == 0 {
		return "", errors.New("no remote branch to detect project name")
	}

	remoteURLBytes, err := repo.output("remote", "get-url", strings.Trim(string(remoteNameBytes), "\n"))
	if err != nil {
		return "", errors.WithStack(err)
	}

	submatchParts := regexp.MustCompile(`/(.+?).git`).FindStringSubmatch(strings.Trim(string(remoteURLBytes), "\n"))
	if len(submatchParts) != 2 {
		return "", errors.New("unknown format of `remote get-url`")
	}

	return submatchParts[1], nil
}

func (repo *repoManager) run(args ...string) error {
	return repo.executor.RunWithWorkDir(repo.repoPath, args...)
}

func (repo *repoManager) output(args ...string) ([]byte, error) {
	output, err := repo.executor.OutputWithWorkDir(repo.repoPath, args...)
	return output, err
}

func (repo repoManager) runWithOpts(args []string, opts ...executor.Opt) error {
	opts = append(opts, executor.WithWorkdir(repo.repoPath))
	return repo.executor.RunWithOpts(args, opts...)
}
