package interfaces

type Client interface {
	Handshake() error
	FirstConnection() error
	QueryUID(uid string) (string, error)
}
