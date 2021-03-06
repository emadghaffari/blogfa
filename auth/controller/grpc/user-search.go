package grpc

import (
	"blogfa/auth/domain/jwt"
	"blogfa/auth/domain/permission"
	"blogfa/auth/domain/provider"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SearchUser method for search user by users fields
func (a *Auth) SearchUser(req *pb.SearchRequest, stream pb.Auth_SearchUserServer) error {
	span := jtrace.Tracer.StartSpan("search-user")
	defer span.Finish()
	span.SetTag("service", "start to search")

	// verify the jwt token
	if _, err := jwt.Model.Verify(req.GetToken()); err != nil {
		return status.Errorf(codes.Internal, "user not verified: %s", err.Error())
	}

	// convert string to int
	from, err := strconv.Atoi(req.GetFrom())
	if err != nil {
		return status.Errorf(codes.Internal, "invalid from number")
	}

	// convert string to int
	to, err := strconv.Atoi(req.GetTo())
	if err != nil {
		return status.Errorf(codes.Internal, "invalid to number")
	}

	users, err := grpc.Service.SearchUser(jtrace.Tracer.ContextWithSpan(context.Background(), span), from, to, req.GetSearch())
	if err != nil {
		return err
	}

	// provider list need to change
	for _, user := range users {
		err := stream.Send(&pb.User{
			Username:  user.Username,
			Gender:    pb.User_Gender(pb.User_Gender_value[user.Gender]),
			Name:      user.Name,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
			BirthDate: user.BirthDate,
			Role: &pb.Role{
				Name:        user.Role.Name,
				Permissions: permission.ToList(user.Role.Permissions),
			},
			Providers: provider.Model.ToProto(user.Provider),
		})
		if err != nil {
			return status.Errorf(codes.Internal, "internal error for get user")
		}
	}

	return nil
}
