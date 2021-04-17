package service

import (
	"blogfa/auth/model/jwt"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"fmt"
)

// SearchUser method for search user by users fields
func (a *Auth) SearchUser(req *pb.SearchRequest, stream pb.Auth_SearchUserServer) error {

	span := jtrace.Tracer.StartSpan("search-user")
	defer span.Finish()
	span.SetTag("service", "start to search")

	// verify the jwt token
	if _, err := jwt.Model.Verify(req.GetToken()); err != nil {
		return fmt.Errorf("user not verified: %s", err.Error())
	}

	return nil
}
