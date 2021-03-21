package config

type nats struct {
	Username  string   `yaml:"nats.username"`
	Password  string   `yaml:"nats.password"`
	Auth      bool     `yaml:"nats.auth"`
	Endpoints []string `yaml:"nats.endpoints"`
}
