package v1

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
	"userCRUD/internal/common/constants"
	"userCRUD/internal/common/deps"
	"userCRUD/internal/user/infrastructure/persistence"
)

const (
	AuthHeader  = "authorization"
	BasicPrefix = "Basic "
)

var (
	ErrNoAuthHeader           = errors.New("no authorization header provided")
	ErrNoMetadata             = errors.New("no metadata provided")
	ErrNoBasicHeader          = errors.New("no basic header")
	ErrInvalidBasicAuthFormat = errors.New("invalid basic auth format")
)

type basicAuthCreds struct {
	username string
	password string
}

func TraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	traceID := uuid.New().String()
	ctx = context.WithValue(ctx, constants.TraceId, traceID)

	return handler(ctx, req)
}

func NewAuthInterceptor(ur persistence.UserRepository, l deps.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		creds, err := getCredsFromHeader(ctx)
		if err != nil {
			return handler(ctx, req)
		}

		user, err := ur.GetUserByUsernameAndPassword(ctx, creds.username, creds.password)
		if err != nil {
			return handler(ctx, req)
		}
		if user != nil {
			newCtx := context.WithValue(ctx, constants.UserContextKey, user)
			return handler(newCtx, req)
		}

		return handler(ctx, req)
	}
}

func getCredsFromHeader(ctx context.Context) (*basicAuthCreds, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrNoMetadata
	}

	authHeaders, ok := md[AuthHeader]
	if !ok || len(authHeaders) == 0 {
		return nil, ErrNoAuthHeader
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, BasicPrefix) {
		return nil, ErrNoBasicHeader
	}

	creds, err := decodeBasicAuth(authHeader)
	if err != nil {
		return nil, ErrInvalidBasicAuthFormat
	}

	return creds, nil
}

func decodeBasicAuth(authHeader string) (*basicAuthCreds, error) {
	authBase64 := strings.TrimPrefix(authHeader, BasicPrefix)
	authBytes, err := base64.StdEncoding.DecodeString(authBase64)
	if err != nil {
		return nil, err
	}

	parts := strings.SplitN(string(authBytes), ":", 2)
	if len(parts) != 2 {
		return nil, ErrInvalidBasicAuthFormat
	}

	return &basicAuthCreds{
		username: parts[0],
		password: parts[1],
	}, nil
}
