package v1

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "userCRUD/api/proto"
	"userCRUD/internal/common"
	"userCRUD/internal/common/deps"
	"userCRUD/internal/user/domain/command"
	"userCRUD/internal/user/domain/model"
	"userCRUD/internal/user/infrastructure/persistence"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	l  deps.Logger
	uc *command.User
}

func NewServer(l deps.Logger, uc *command.User) *Server {
	return &Server{
		l:  l,
		uc: uc,
	}
}

func (s *Server) NewUser(ctx context.Context, req *pb.NewUserRequest) (*pb.UserResponse, error) {
	user, err := s.uc.CreateUser(ctx, &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Admin:    req.Admin,
	})

	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.UserResponse{
		Id:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Admin:    user.Admin,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	user, err := s.uc.UpdateUser(ctx, &model.UpdateUser{
		ID:       req.Id,
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
		Admin:    req.Admin,
	})

	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.UserResponse{
		Id:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Admin:    user.Admin,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.uc.DeleteUser(ctx, &model.UserByID{
		ID: req.Id,
	})

	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.DeleteUserResponse{}, nil
}

func (s *Server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	user, err := s.uc.GetUserByID(ctx, &model.UserByID{
		ID: req.Id,
	})

	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.UserResponse{
		Id:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Admin:    user.Admin,
	}, nil
}

func (s *Server) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.UserResponse, error) {
	user, err := s.uc.GetUserByUsername(ctx, &model.UserByUsername{
		Username: req.Username,
	})

	if err != nil {
		return nil, handleGRPCError(err)
	}

	return &pb.UserResponse{
		Id:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		Admin:    user.Admin,
	}, nil
}

func (s *Server) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := s.uc.GetUsers(ctx, &common.Pagination{
		Page:     req.Page,
		PageSize: req.PageSize,
	})

	if err != nil {
		return nil, handleGRPCError(err)
	}

	usersResp := make([]*pb.UserResponse, len(users))
	for i, u := range users {
		usersResp[i] = &pb.UserResponse{
			Id:       u.ID,
			Email:    u.Email,
			Username: u.Username,
			Admin:    u.Admin,
		}
	}

	return &pb.GetUsersResponse{
		Users: usersResp,
	}, nil
}

func handleGRPCError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, persistence.ErrUserNotFound):
		return status.Errorf(codes.NotFound, err.Error())
	case errors.Is(err, command.ErrNotEnoughPermissions), errors.Is(err, command.ErrAuthFailed):
		return status.Errorf(codes.PermissionDenied, err.Error())
	default:
		return status.Errorf(codes.InvalidArgument, err.Error())
	}
}
