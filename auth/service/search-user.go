package service

import pb "blogfa/auth/proto"

// SearchUser method for search user by users fields
func (a *Auth) SearchUser(req *pb.SearchRequest, stream pb.Auth_SearchUserServer) error {
	return nil
}
