package net

import (
	"context"
	gonet "net"
	"syscall"

	"github.com/godzie44/go-uring/reactor"
	"github.com/godzie44/go-uring/uring"
)

type dialer struct {
	reactor *reactor.NetworkReactor
}

func (d *dialer) DialContext(ctx context.Context, network, address string) (gonet.Conn, error) {
	sockFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, 0)
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
