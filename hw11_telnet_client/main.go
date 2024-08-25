package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connection timeout")
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Not enough args")
		os.Exit(-1)
	}

	address, err := validateHostPort(flag.Args()[0], flag.Args()[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid args: %s\n", err)
		os.Exit(-1)
	}

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	err = client.Connect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Faild to establish connection to server: %s\n", err)
		os.Exit(0)
	}
	fmt.Fprintf(os.Stderr, "Establish connection to %s\n", address)

	go func() {
		for {
			err := client.Send()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Programm terminate: %s\n", err)
				err = client.Close()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Connection close err: %s:", err)
				}
				os.Exit(0)
			}
		}
	}()

	go func() {
		for {
			client.Receive()
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}()

	<-ctx.Done()
}
