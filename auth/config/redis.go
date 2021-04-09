package config

// redis struct
type redis struct {
	Username string `yaml:"redis.username"`
	Password string `yaml:"redis.password"`
	DB       int    `yaml:"redis.db"`
	Host     string `yaml:"redis.host"`
	Logger   bool   `yaml:"redis.logger"`
}
