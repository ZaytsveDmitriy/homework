package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrConnectionClosed       = errors.New("connection closed")
	ErrConectionCloseRequest  = errors.New("connection close user request")
	ErrUnexpected             = errors.New("unexpected error")
	ErrServerConnectionClosed = errors.New("server closed connection unexpectedly")
	ErrServerTimeOut          = errors.New("server time out")
)

type SimpleTelnetClient struct {
	connect   net.Conn
	netReader *bufio.Reader
	address   string
	timeout   time.Duration
	in        *bufio.Reader
	out       io.Writer
}

func (c *SimpleTelnetClient) Connect() error {
	var err error
	c.connect, err = net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.netReader = bufio.NewReader(c.connect)
	return nil
}

func (c *SimpleTelnetClient) Close() error {
	err := c.connect.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *SimpleTelnetClient) Send() error {
	const op = "client.Send"
	d, err := c.in.ReadBytes('\n')
	if err != nil && errors.Is(err, io.EOF) {
		return fmt.Errorf("%s: %w", op, ErrConectionCloseRequest)
	}

	_, err = c.connect.Write(d)
	if err != nil {
		switch {
		case errors.Is(err, net.ErrClosed):
			return fmt.Errorf("%s: %w", op, ErrServerConnectionClosed)
		case os.IsTimeout(err):
			return fmt.Errorf("%s: %w", op, ErrServerTimeOut)
		default:
			return fmt.Errorf("%s: %w", op, ErrUnexpected)
		}
	}
	return nil
}

func (c *SimpleTelnetClient) Receive() error {
	const op = "client.Receive"
	d, err := c.netReader.ReadString('\n')
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			return fmt.Errorf("%s: %w", op, ErrConnectionClosed)
		default:
			return fmt.Errorf("%s: %w", op, ErrUnexpected)
		}
	}
	_, err = c.out.Write([]byte(d))
	if err != nil {
		return fmt.Errorf("unexpected error")
	}
	return nil
}

var (
	ErrHostInvalid = errors.New("invalid host")
	ErrPortInvalid = errors.New("invalid port")
)

func validateHostPort(host, port string) (address string, err error) {
	const op = "validateHostPort"
	const localhost = "localhost"

	// TODO: more adeqate validator

	if !strings.Contains(host, ".") && host != localhost {
		return "", fmt.Errorf("%s: %w", op, ErrHostInvalid)
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 0 || portNum > 65535 {
		return "", fmt.Errorf("%s: %w", op, ErrPortInvalid)
	}

	return net.JoinHostPort(host, port), nil
}
