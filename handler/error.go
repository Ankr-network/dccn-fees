package handler

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"runtime"
)

var (
	ErrAccessDenied = errors.New("access denied")
)

func handleError(err error) error {
	_, f, l, _ := runtime.Caller(1)
	log.Printf("file=%s,line=%d,msg=%v", f, l, err)

	// convert to grpc error
	switch {
	case errors.Is(err, ErrAccessDenied):
		return status.Errorf(codes.PermissionDenied, "AccessDenied:%v", err.Error())
	default:
		return status.Errorf(codes.Unknown, "InternalError:%v", err.Error())
	}
}
