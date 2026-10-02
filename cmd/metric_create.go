package main

import (
	"errors"
	"flag"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/mem"
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
		return errors.New("Name could not be empty, use -name flag to pass something")
	}

	// NOTE: here i raise the same question: why the hell should i keep the id counter somewhere
	// outside of the actual metric creating place
	// anyways, here we are creating the metric
	// NOTE: answer
	// we are passing the pointer inside of the fucntion the function is filling it. that's the job basically
	// we are not assigning ID or other important fields outside, that's the job of the inside function.
	// we are just working on the same object inside & outside for convenience
	met := meth.Metric{
		Name: *name,
	}
	// here we create inmemory metric service. essentially an abstract controller
	ms := mem.NewMetricService()
	// creating metric and storing it in memory
	err := ms.CreateMetric(&met)

	return err
}
