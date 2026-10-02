package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	// "github.com/demolemo/meth/mem"
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	err := Run(os.Args[1:])
	if err != nil {
		fmt.Printf(err.Error())
	} else {
		fmt.Printf("Current error is no error")
	}

	// <-ctx.Done()
}

func Run(args []string) error {
	var cmd string
	if len(args) > 1 {
		cmd, args = args[0], args[1:]
	}
	// for debugging and priting commands
	// fmt.Printf("cmd var: %s\n", cmd)
	// for _, arg := range args {
	// 	fmt.Printf("arg var: %s\n", arg)
	// }
	switch cmd {
	// NOTE: I see why we need the swtich, because it's related to different modes we want to execute the program in
	// first we have a direct execution which is dispatched inside of metric.go
	// second we have honest usage which we write here
	// thirdly we have something for debug? i don't remember what the freaking dial.go did
	case "create":
		return (&MetricCreateCommand{}).Run(args)
	case "delete":
		return (&MetricDeleteCommand{}).Run(args)
	}
	return nil
}
