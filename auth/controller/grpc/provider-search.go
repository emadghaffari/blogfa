package grpc

import (
	"blogfa/auth/domain/jwt"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"blogfa/auth/service/grpc"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SearchProvider method for search user by provider fields
func (a *Auth) SearchProvider(req *pb.SearchRequest, stream pb.Auth_SearchProviderServer) error {
	span := jtrace.Tracer.StartSpan("search-provider")
	defer span.Finish()
	span.SetTag("service", "start to search")

	// verify the jwt token
	if _, err := jwt.Model.Verify(req.GetToken()); err != nil {
		return status.Errorf(codes.Internal, "provider not verified: %s", err.Error())
	}

	providers, err := grpc.Service.SearchProvider(jtrace.Tracer.ContextWithSpan(stream.Context(), span), req)
	if err != nil {
		return err
	}

	for _, provider := range providers {
		stream.Send(&pb.Provider{
			ID:          strconv.Itoa(int(provider.ID)),
			FixedNumber: provider.FixedNumber,
			Company:     provider.Company,
			Card:        provider.Card,
			CardNumber:  provider.CardNumber,
			ShebaNumber: provider.ShebaNumber,
			Address:     provider.Address,
			User: &pb.User{
				Username:  provider.User.Username,
				Name:      provider.User.Name,
				LastName:  provider.User.LastName,
				Gender:    pb.User_Gender(pb.User_Gender_value[provider.User.Gender]),
				Phone:     provider.User.Phone,
				Email:     provider.User.Email,
				BirthDate: provider.User.BirthDate,
			},
			Token: req.GetToken(),
		})
		if err != nil {
			return status.Errorf(codes.Internal, "internal error for get provider")
		}
	}

	return nil
}
