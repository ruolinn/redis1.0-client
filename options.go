package redis

import (
	"net"
	"time"
)

type Options struct {
	NetWork     string
	Addr        string
	Dialer      func() (net.Conn, error)
	DB          int
	MaxRetries  int
	DialTimeout time.Duration
}

func (opt *Options) init() {
	if opt.NetWork == "" {
		opt.NetWork = "tcp"
	}

	if opt.Dialer == nil {
		opt.Dialer = func() (net.Conn, error) {
			conn, err := net.DialTimeout(opt.NetWork, opt.Addr, opt.DialTimeout)

			return conn, err
		}
	}

	if opt.DialTimeout == 0 {
		opt.DialTimeout = 5 * time.Second
	}
}
