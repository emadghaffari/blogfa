package config

type consul struct {
	Address       string
	Scheme        string
	Tags          []string
	CheckHttp     string
	CheckInterval string
}
