package service

import (
	"blogfa/auth/model/jwt"
	"blogfa/auth/model/user"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
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

	// search users
	rows, err := user.Model.Search(jtrace.Tracer.ContextWithSpan(context.Background(), span), from, to, req.GetSearch())
	if err != nil {
		return status.Errorf(codes.Internal, "internal error for search users")
	}
	defer rows.Close()

	for rows.Next() {
		usr := user.User{}
		if err := rows.Scan(&usr.Username, &usr.Name, &usr.LastName, &usr.Phone, &usr.Email, &usr.Gender, &usr.RoleID); err != nil {
			status.Errorf(codes.Internal, fmt.Sprintf("internal error for get users"))
		}

		stream.Send(&pb.User{
			Gender:    pb.User_Gender(pb.User_Gender_value[usr.Gender]),
			Name:      usr.Name,
			LastName:  usr.LastName,
			Phone:     usr.Phone,
			Email:     usr.Email,
			BirthDate: usr.BirthDate,
			Role: &pb.Role{
				Name: usr.Role.Name,
			},
		})
		// if err != nil {
		// 	status.Errorf(codes.Internal, fmt.Sprintf("internal error for get user"))
		// }
	}

	fmt.Println("DONE")

	return nil
}
