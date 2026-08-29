package main

import (
	"flag"
	"fmt"
	"strings"
)

// MetricCommand represents a collection of metric-related subcommands.
type MetricCommand struct{}

// Run executes the command which delegates to other subcommands.
// NOTE: this was the false part of the flow, we are currently not using this bs
// and going directly to the subcommand that we need
func (c *MetricCommand) Run(args []string) error {
	// Shift off the subcommand name, if available.
	var cmd string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		cmd, args = args[0], args[1:]
	}

	// Delegete to the appropriate subcommand.
	switch cmd {
	case "", "list":
		return (&MetricListCommand{}).Run(args)
	case "create":
		return (&MetricCreateCommand{}).Run(args)
	case "delete":
		return (&MetricDeleteCommand{}).Run(args)
	case "help":
		c.usage()
		return flag.ErrHelp
	default:
		return fmt.Errorf("wtf dial %s: unknown command", cmd)
	}
}

// usage prints the subcommand usage to STDOUT.
// i don't like that the "metric" will be used everywhere in that command definition
// i think that this is redundant
func (c *MetricCommand) usage() {
	fmt.Println(`
Manage metrics you own or are a member of.

Usage:

	meth metric <command> [arguments]

The commands are:

	list        list all available metrics
	create      create a new metric
	delete      remove an existing metric
`[1:])
}
