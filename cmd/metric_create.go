package main

import "flag"

type MetricCreator struct{}

// first version of the command, create without the flags
// name - if missing we fail
func Run(args []string) error {
	fs := flag.NewFlagSet("metric-create", flag.ContinueOnError)
	name := fs.String("name", "", "Give this metric a name will ya")
	metricService := NewDialService()

	_ = name
	return nil
}
