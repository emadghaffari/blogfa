package config

var (
	// Global config
	Global GlobalConfig
)

// GlobalConfig is base of configs we need for project
type GlobalConfig struct {
	Environment string `yaml:"environment"`
	Service     service
	Jaeger      jaeger
	Log         loggingConfig
	ETCD        etcd
	Redis       redis
	MYSQL       database
}

// GetService is method for get a service struct with default vaules by config file
func (g *GlobalConfig) GetService() interface{} {
	service := struct {
		Name string
		GRPC struct {
			Port string
			Host string
		}
	}{
		Name: Global.Service.Name,
		GRPC: struct {
			Port string
			Host string
		}{
			Port: Global.Service.GRPC.Port,
			Host: Global.Service.GRPC.Host,
		},
	}

	return service
}

// Service details
type service struct {
	Name string `yaml:"service.name"`
	GRPC struct {
		Host     string `yaml:"grpc.host"`
		Port     string `yaml:"grpc.port"`
		Endpoint string `yaml:"grpc.endpoint"`
	}
	HTTP struct {
		Host     string `yaml:"http.host"`
		Port     string `yaml:"http.port"`
		Endpoint string `yaml:"http.endpoint"`
	}
	DEBUG struct {
		Host     string `yaml:"debug.host"`
		Port     string `yaml:"debug.port"`
		Endpoint string `yaml:"debug.endpoint"`
	}
}

// Jaeger tracer
type jaeger struct {
	HostPort string `yaml:"jaeger.hostPort"`
	LogSpans bool   `yaml:"jaeger.logSpans"`
}

// LoggingConfig struct
type loggingConfig struct {
	DisableColors    bool `json:"disable_colors" yaml:"log.disableColors"`
	QuoteEmptyFields bool `json:"quote_empty_fields" yaml:"log.quoteEmptyFields"`
}

type etcd struct {
	Endpoints []string `json:"endpoints" yaml:"etcd.endpoints"`
	WatchList []string `json:"watch_list" yaml:"etcd.watchList"`
}

// redis struct
type redis struct {
	Address string `json:"address" yaml:"redis.address"`
}

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
