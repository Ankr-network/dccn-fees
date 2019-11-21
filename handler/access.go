package handler

import (
	"context"
	"fmt"
	"github.com/Ankr-network/dccn-common/access"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func newResource(teamID, path string) string {
	return fmt.Sprintf("ankr:fee:%s:%s", teamID, path)
}

func checkAccess(ctx context.Context, uid, res, action string) error {
	log.Printf("uid=%s,res=%s,action=%s", uid, res, action)
	ok, err := access.Authorize(ctx, uid, res, action)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	if !ok {
		return ErrAccessDenied
	}
	return nil
}
