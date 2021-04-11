package jwt

import (
	"context"

	"go.uber.org/zap"
)

var (
	// JWT variable instance of intef
	Model  intef = &jwt{}
	logger *zap.Logger
)

// jwt meths interface
type intef interface {
	Generate(ctx context.Context, model interface{}) (*jwt, error)
	GenerateJWT() (*jwt, error)
	genRefJWT(td *jwt) error
	store(ctx context.Context, model interface{}, td *jwt) error
	Get(ctx context.Context, token string, response interface{}) error
}

// jwt struct
type jwt struct {
	AccessToken  string `json:"at"`
	RefreshToken string `json:"rt"`
	AccessUUID   string `json:"uuid"`
	RefreshUUID  string `json:"rau"`
	AtExpires    int64  `json:"exp"`
	RtExpires    int64  `json:"rexp"`
}
