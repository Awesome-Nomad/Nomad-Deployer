package config

type DirectConnection struct {
	address string
}

func (c *DirectConnection) Init(address string) error {
	c.address = address
	return nil
}

func (c *DirectConnection) Destroy() error {
	return nil
}

func (c *DirectConnection) GetAddress() (string, error) {
	return HTTPURLEnhancer(c.address), nil

}
