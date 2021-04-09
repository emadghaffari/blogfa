package config

import "time"

// redis struct
type redis struct {
	Username            string        `yaml:"redis.username"`
	Password            string        `yaml:"redis.password"`
	DB                  int           `yaml:"redis.db"`
	Host                string        `yaml:"redis.host"`
	Logger              bool          `yaml:"redis.logger"`
	SMSDuration         time.Duration `yaml:"redis.smsDuration"`
	SMSCodeVerification time.Duration `yaml:"redis.smsCodeVerification"`
	UserDuration        time.Duration `yaml:"redis.userDuration"`
}
