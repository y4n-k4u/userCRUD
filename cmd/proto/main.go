package main

import (
	"context"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	pb "userCRUD/api/proto"
	"userCRUD/internal/common/deps"
	"userCRUD/internal/user/domain/command"
	"userCRUD/internal/user/infrastructure/persistence"
	"userCRUD/internal/user/infrastructure/transport/proto/v1"
)

func main() {
	container := buildContainer()

	if err := container.Invoke(runApp); err != nil {
		log.Printf("Critical error: %v\n", err)
		os.Exit(1)
	}
}

func buildContainer() *dig.Container {
	container := dig.New()

	container.Provide(deps.NewZapLogger, dig.As(new(deps.Logger)))
	container.Provide(deps.NewGoPlaygroundValidator, dig.As(new(deps.Validator)))
	container.Provide(persistence.NewUserRepositoryMemory, dig.As(new(persistence.UserRepository)))
	container.Provide(func(ur persistence.UserRepository, v deps.Validator) *command.User {
		return command.NewUserCommand(ur, v)
	})

	container.Provide(func(l deps.Logger, uc *command.User, ur persistence.UserRepository) *grpc.Server {
		return newGRPCServer(uc, ur, l)
	})

	return container
}

func newGRPCServer(uc *command.User, ur persistence.UserRepository, l deps.Logger) *grpc.Server {
	ai := v1.NewAuthInterceptor(ur, l)
	chain := grpc.ChainUnaryInterceptor(
		v1.TraceInterceptor,
		ai,
	)
	server := grpc.NewServer(chain)
	pb.RegisterUserServiceServer(server, v1.NewServer(l, uc))

	return server
}

func runApp(logger deps.Logger, s *grpc.Server) {
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		logger.Error(context.Background(), "failed to listen", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Error(context.Background(), "failed to serve", "error", err)
		}
	}()
	logger.Info(context.Background(), "Server started on :50051")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	logger.Info(context.Background(), "Shutting down server...")

	s.GracefulStop()
	logger.Info(context.Background(), "Server stopped")
}
