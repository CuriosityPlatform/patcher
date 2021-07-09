package os

import "patcher/pkg/common/infrastructure/executor"

type HostInfoProvider interface {
	Username() (string, error)
	DeviceName() (string, error)
}

func NewHostInfoProvider() HostInfoProvider {
	return &hostInfoProvider{}
}

type hostInfoProvider struct {
}

func (provider *hostInfoProvider) Username() (string, error) {
	exec, err := executor.New("whoami")
	if err != nil {
		return "", err
	}

	usernameBytes, err := exec.Output()
	if err != nil {
		return "", err
	}

	return string(usernameBytes), nil
}

func (provider *hostInfoProvider) DeviceName() (string, error) {
	exec, err := executor.New("hostname")
	if err != nil {
		return "", err
	}

	hostnameBytes, err := exec.Output()
	if err != nil {
		return "", err
	}

	return string(hostnameBytes), nil
}
