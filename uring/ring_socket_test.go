package uring

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"syscall"
	"testing"
)

func TestSocket(t *testing.T) {
	ring, err := New(64)
	require.NoError(t, err)
	defer ring.Close()

	domain := syscall.AF_INET
	typ := syscall.SOCK_STREAM
	protocol := syscall.IPPROTO_TCP

	err = ring.QueueSQE(Socket(domain, typ, protocol), 0, 0)
	require.NoError(t, err)
	_, err = ring.Submit()
	require.NoError(t, err)

	cqe, err := ring.WaitCQEvents(1)
	require.NoError(t, err)
	require.NoError(t, cqe.Error())

	socketFd := cqe.Res
	require.Positive(t, socketFd)
	fmt.Printf("Created socket with fd: %d\n", socketFd)

	ring.SeenCQE(cqe)

	defer syscall.Close(int(socketFd))
}
