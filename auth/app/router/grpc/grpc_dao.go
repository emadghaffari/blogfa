package grpc

var (
	Router auth = &Auth{}
)

// Auth service
type Auth struct{}

type auth interface {
}
