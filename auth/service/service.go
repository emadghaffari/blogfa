package service

import (
	"blogfa/auth/model/provider"
	"blogfa/auth/model/user"
	"blogfa/auth/pkg/cript"
	"blogfa/auth/pkg/jtrace"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
)

// Auth service
type Auth struct{}

// RegisterUser, for create a new user
func (a *Auth) RegisterUser(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	span := jtrace.Tracer.StartSpan("register-user")
	defer span.Finish()
	span.SetTag("register", "register user")

	password, err := cript.Hash(req.GetPassword())
	if err != nil {
		return &pb.UserRegisterResponse{Message: "ERROR"}, fmt.Errorf("error in hash password: %s", err.Error())
	}

	// create new user requested.
	_, err = user.Model.Register(jtrace.Tracer.ContextWithSpan(ctx, span), user.User{
		Username:  req.GetUsername(),
		Password:  &password,
		Name:      req.GetName(),
		LastName:  req.GetLastName(),
		Phone:     req.GetPhone(),
		Email:     req.GetEmail(),
		BirthDate: req.GetBirthDate(),
		Gender:    req.GetGender().String(),
		RoleID:    1, // USER
	})
	if err != nil {
		return &pb.UserRegisterResponse{Message: "ERROR"}, fmt.Errorf("error in store user: %s", err.Error())
	}

	child := jtrace.Tracer.ChildOf(span, "register")
	child.SetTag("register", "after register user")
	defer child.Finish()

	return &pb.UserRegisterResponse{Message: "DONE"}, nil
}

// RegisterProvider, for create new provider
func (a *Auth) RegisterProvider(ctx context.Context, req *pb.ProviderRegisterRequest) (*pb.ProviderRegisterResponse, error) {
	span := jtrace.Tracer.StartSpan("register-provider")
	defer span.Finish()
	span.SetTag("register", "register provider")

	password, err := cript.Hash(req.GetPassword())
	if err != nil {
		return &pb.ProviderRegisterResponse{Message: "ERROR"}, fmt.Errorf("error in hash password: %s", err.Error())
	}

	// create new user requested.
	user, err := user.Model.Register(jtrace.Tracer.ContextWithSpan(ctx, span), user.User{
		Username:  req.GetUsername(),
		Password:  &password,
		Name:      req.GetName(),
		LastName:  req.GetLastName(),
		Phone:     req.GetPhone(),
		Email:     req.GetEmail(),
		BirthDate: req.GetBirthDate(),
		Gender:    req.GetGender().String(),
		RoleID:    1, // USER
	})
	if err != nil {
		return &pb.ProviderRegisterResponse{Message: "ERROR"}, fmt.Errorf("error in store user: %s", err.Error())
	}

	if err := provider.Model.Register(jtrace.Tracer.ContextWithSpan(ctx, span), provider.Provider{
		UserID:      user.ID,
		FixedNumber: req.GetFixedNumber(),
		Company:     req.GetCompany(),
		Card:        req.GetCard(),
		CardNumber:  req.GetCardNumber(),
		ShebaNumber: req.GetShebaNumber(),
		Address:     req.GetAddress(),
	}); err != nil {
		return &pb.ProviderRegisterResponse{Message: "ERROR"}, fmt.Errorf("error in store new provider: %s", err.Error())
	}

	child := jtrace.Tracer.ChildOf(span, "register")
	child.SetTag("register", "after register provider")
	defer child.Finish()

	return &pb.ProviderRegisterResponse{}, nil
}
func (a *Auth) UPLogin(ctx context.Context, req *pb.UPLoginRequest) (*pb.UPLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("up-login")
	defer span.Finish()

	return &pb.UPLoginResponse{}, nil
}
func (a *Auth) PLogin(ctx context.Context, req *pb.PLoginRequest) (*pb.PLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("p-login")
	defer span.Finish()

	return &pb.PLoginResponse{}, nil
}
