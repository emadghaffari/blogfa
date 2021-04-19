package service

import pb "blogfa/auth/proto"

// SearchProvider method for search user by provider fields
func (a *Auth) SearchProvider(req *pb.SearchRequest, stream pb.Auth_SearchProviderServer) error {
	return nil
}
