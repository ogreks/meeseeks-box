package secret

type Crypto interface {
	Name() string
	Sign(msg string, secret string) ([]byte, error)
}
