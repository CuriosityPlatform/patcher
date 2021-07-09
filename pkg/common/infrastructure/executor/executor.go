//nolint
package executor

import "os/exec"

type Executor interface {
	Output(args ...string) ([]byte, error)
	OutputWithWorkDir(workDir string, args ...string) ([]byte, error)
	Run(args ...string) error
	RunWithWorkDir(workDir string, args ...string) error
}

func New(executable string) (Executor, error) {
	_, err := exec.LookPath(executable)
	if err != nil {
		return nil, err
	}

	return &executor{executable: executable}, nil
}

type executor struct {
	executable string
}

func (e *executor) Output(args ...string) ([]byte, error) {
	cmd := exec.Command(e.executable, args...)
	return cmd.Output()
}

func (e *executor) OutputWithWorkDir(workDir string, args ...string) ([]byte, error) {
	cmd := exec.Command(e.executable, args...)
	cmd.Dir = workDir
	return cmd.Output()
}

func (e *executor) Run(args ...string) error {
	cmd := exec.Command(e.executable, args...)
	return cmd.Run()
}

func (e *executor) RunWithWorkDir(workDir string, args ...string) error {
	cmd := exec.Command(e.executable, args...)
	cmd.Dir = workDir
	return cmd.Run()
}
