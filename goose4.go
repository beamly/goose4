package goose4

// Goose4 holds goose4 configuration and provides functions thereon
type Goose4 struct {
	c Config
}

// NewGoose4 returns a Goose4 object to be used as net/http handler
func NewGoose4(c Config) (g Goose4, err error) {
	g.c = c

	return
}
