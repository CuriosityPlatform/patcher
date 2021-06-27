package transport

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"patcher/pkg/patcherservice/app"
)

func translateError(err error) error {
	switch errors.Cause(err) {
	case app.ErrPatchAlreadyApplied:
	case app.ErrPatchCantAddPatchWitEmptyContent:
		return status.Error(codes.InvalidArgument, err.Error())
	case app.ErrPatchNotFound:
		return status.Error(codes.NotFound, err.Error())
	}

	return err
}
