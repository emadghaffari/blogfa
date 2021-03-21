package config

// database struct
type database struct {
	Username    string `yaml:"mysql.username"`
	Password    string `yaml:"mysql.password"`
	Host        string `yaml:"mysql.host"`
	Schema      string `yaml:"mysql.schema"`
	Driver      string `yaml:"mysql.driver"`
	Automigrate bool   `yaml:"mysql.automigrate"`
	Logger      bool   `yaml:"mysql.logger"`
	Namespace   string
}
