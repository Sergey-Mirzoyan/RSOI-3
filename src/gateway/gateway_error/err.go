package gateway_error

type GatewayError struct {
	Err error
	Code int
}

const (
	Ok = 0
	Internal = 1
	User = 2
)
