package service

import (
	"blogfa/auth/broker"
	"blogfa/auth/config"
	"blogfa/auth/model/jwt"
	"blogfa/auth/model/user"
	"blogfa/auth/pkg/jtrace"
	"blogfa/auth/pkg/token"
	pb "blogfa/auth/proto"
	"context"
	"fmt"
	"net/http"
)

// Auth service
type Auth struct{}

// PLogin, login user with phone number with sms code
func (a *Auth) PLogin(ctx context.Context, req *pb.PLoginRequest) (*pb.PLoginResponse, error) {
	span := jtrace.Tracer.StartSpan("p-login")
	defer span.Finish()
	span.SetTag("login", "phone login")

	// get user with email or username
	user, err := user.Model.Get(jtrace.Tracer.ContextWithSpan(ctx, span), "users", "phone = ? ", req.GetPhone())
	if err != nil {
		return &pb.PLoginResponse{
			Message: "invalid phone number",
			Status: &pb.Status{
				Code:    http.StatusInternalServerError,
				Message: "invalid phone number",
			},
		}, fmt.Errorf("invalid phone number")
	}

	// generate jwt token
	jwt, err := jwt.Model.GenerateJWT()
	if err != nil {
		return &pb.PLoginResponse{
			Message: "error in generate accessToken try after 10 seconds!",
			Status: &pb.Status{
				Code:    http.StatusInternalServerError,
				Message: "error in generate accessToken try after 10 seconds!",
			},
		}, fmt.Errorf("error in generate accessToken try after 10 seconds!")
	}

	// make a map for jwt and user
	data := make(map[string]interface{}, 2)
	data["jwt"] = *jwt
	data["user"] = user
	notif := config.SMS{
		Service: config.Global.Service.Name,
		Token:   token.Generate(25),
		Data:    data,
		To:      user.Phone,
	}

	// publish to nats channel
	broker.Nats.Publish(ctx, "service.notification.sms.auth", notif)

	// return response for check the phone
	return &pb.PLoginResponse{
		Message: "check your phone!",
		Token:   jwt.AccessUUID,
		Status: &pb.Status{
			Code:    http.StatusOK,
			Message: "successfully",
		},
	}, nil
}
