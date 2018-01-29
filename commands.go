package redis

type cmdable struct {
	process func(cmd Cmder) error
}

func (c *cmdable) setProcessor(fn func(Cmder) error) {
	c.process = fn
}

func (c *cmdable) Ping() *StatusCmd {
	cmd := NewStatusCmd("ping")
	c.process(cmd)

	return cmd
}
