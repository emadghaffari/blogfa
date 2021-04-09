package jwt

import (
	"blogfa/auth/config"
	"blogfa/auth/database/redis"
	"blogfa/auth/pkg/token"
	"context"
	"encoding/json"
	"fmt"
	"time"

	jjwt "github.com/dgrijalva/jwt-go"
)

var (
	// JWT variable instance of intef
	JWT intef = &jwt{}
)

// jwt meths interface
type intef interface {
	Generate(ctx context.Context, model interface{}) (*jwt, error)
	genJWT() (*jwt, error)
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

// Generate new jwt token and store into redis DB
func (j *jwt) Generate(ctx context.Context, model interface{}) (*jwt, error) {

	td, err := j.genJWT()
	if err != nil {
		return nil, err
	}

	if err := j.genRefJWT(td); err != nil {
		return nil, err
	}

	if err := j.store(ctx, model, td); err != nil {
		return nil, err
	}

	return td, nil
}

// generate JWT tokens
func (j *jwt) genJWT() (*jwt, error) {
	// create new jwt
	td := &jwt{}
	td.AtExpires = time.Now().Add(time.Duration(config.Global.Redis.UserDuration)).Unix()
	td.RtExpires = time.Now().Add(time.Duration(config.Global.Redis.UserDuration)).Unix()
	td.AccessUUID = token.Generate(30)
	td.RefreshUUID = token.Generate(60)

	// New MapClaims for access token
	atClaims := jjwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["uuid"] = td.AccessUUID
	atClaims["exp"] = td.AtExpires
	at := jjwt.NewWithClaims(jjwt.SigningMethodHS256, atClaims)

	var err error
	td.AccessToken, err = at.SignedString([]byte(config.Global.JWT.Secret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

// generate refresh tokens
func (j *jwt) genRefJWT(td *jwt) error {
	// New MapClaims for refresh access token
	rtClaims := jjwt.MapClaims{}
	rtClaims["uuid"] = td.RefreshUUID
	rtClaims["exp"] = td.RtExpires
	rt := jjwt.NewWithClaims(jjwt.SigningMethodHS256, rtClaims)

	var err error
	td.RefreshToken, err = rt.SignedString([]byte(config.Global.JWT.RSecret))
	if err != nil {
		return err
	}
	return nil
}

// store into DB
func (j *jwt) store(ctx context.Context, model interface{}, td *jwt) error {
	bt, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("can not marshal data: %s", model)
	}
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	now := time.Now()

	// make map for store in redis
	if err := redis.Storage.Set(ctx, td.AccessUUID, string(bt), at.Sub(now)); err != nil {
		return err
	}
	return nil
}

// Get jwt token from redis
func (j *jwt) Get(ctx context.Context, token string, response interface{}) error {
	if err := redis.Storage.Get(ctx, token, &response); err != nil {
		return err
	}
	return nil
}
