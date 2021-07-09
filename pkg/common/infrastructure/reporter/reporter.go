package reporter

import "github.com/sirupsen/logrus"

type Reporter interface {
	Info(...interface{})
}

func New(isQuiet bool, impl *logrus.Logger) Reporter {
	return &reporter{
		isQuiet: isQuiet,
		impl:    impl,
	}
}

type reporter struct {
	isQuiet bool
	impl    *logrus.Logger
}

func (r *reporter) Info(args ...interface{}) {
	if r.isQuiet {
		return
	}
	r.impl.Info(args...)
}
