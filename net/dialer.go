package net

import (
	"context"
	gonet "net"

	"github.com/lynxai-team/go-uring/reactor"
	"github.com/lynxai-team/go-uring/uring"
	"golang.org/x/sys/unix"
)

type dialer struct {
	reactor *reactor.NetworkReactor
}

func (d *dialer) DialContext(ctx context.Context, network, address string) (gonet.Conn, error) {
	sockFd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}

	if err = setDefaultListenerSockopts(sockFd); err != nil {
		return nil, err
	}

	addr, err := gonet.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	cqeC := make(chan uring.CQEvent)
	connectOp := uring.Connect(uintptr(sockFd), addr)
	d.reactor.Queue(connectOp, func(event uring.CQEvent) {
		cqeC <- event
	})

	cqe := <-cqeC
	if err := cqe.Error(); err != nil {
		return nil, err
	}

	lAddr, err := connectOp.Addr()
	if err != nil {
		return nil, err
	}

	tc := newConn(sockFd, lAddr, addr, d.reactor)
	return tc, nil
}

func NewDialer(reactor *reactor.NetworkReactor) *dialer {
	return &dialer{
		reactor: reactor,
	}
}
