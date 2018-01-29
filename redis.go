package redis

import (
	"fmt"
	"redis/pool"
)

type Client struct {
	opt *Options
	cmdable
}

func NewClient(opt *Options) *Client {
	opt.init()

	client := Client{
		opt: opt,
	}

	client.setProcessor(client.defaultProcess)

	return &client
}

func (c *Client) defaultProcess(cmd Cmder) error {
	cn, err := c.NewConn()

	if err != nil {
		fmt.Println(err)
	}

	if err := writeCmd(cn, cmd); err != nil {
		fmt.Println(err)
	}

	err = cmd.readReply(cn)

	return cmd.Err()
}

func (c *Client) NewConn() (*pool.Conn, error) {
	netConn, err := c.opt.Dialer()

	cn := pool.NewConn(netConn)

	return cn, err
}
