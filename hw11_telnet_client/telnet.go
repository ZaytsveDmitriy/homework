package main

import (
	"bufio"
	"io"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	client := SimpleTelnetClient{
		in:      bufio.NewReaderSize(in, 1500),
		out:     out,
		timeout: timeout,
		address: address,
	}
	return &client
}
