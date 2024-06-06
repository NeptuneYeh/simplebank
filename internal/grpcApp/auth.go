package grpcApp

import (
	"context"
	"fmt"
	"github.com/NeptuneYeh/simplebank/init/auth"
	"github.com/NeptuneYeh/simplebank/tools/token"
	"google.golang.org/grpc/metadata"
	"strings"
)

const (
	authorizationHeader = "authorization"
)

func (c *Module) authUser(ctx context.Context, accessibleRoles []string) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 || strings.ToLower(fields[0]) != "bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	accessToken := fields[1]
	payload, err := auth.MainAuth.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	return payload, nil
}

func hasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}

	return false
}
