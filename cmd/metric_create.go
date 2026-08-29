package main

import (
	"errors"
	"flag"
	"fmt"
)

// first version of the command, create without the flags
// name - if missing we fail
type MetricCreateCommand struct{}

func (c *MetricCreateCommand) Run(args []string) error {
	fs := flag.NewFlagSet("metric-create", flag.ContinueOnError)
	name := fs.String("name", "", "Give this metric a name will ya")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" {
		return errors.New("Name could not be empty")
	}
	return fmt.Errorf("current name is: %s", *name)
}
