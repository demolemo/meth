package main

import (
	"errors"
	"flag"

	"github.com/demolemo/meth"
	"github.com/demolemo/meth/mem"
)

type MetricDeleteCommand struct{}

// stub for an actual future command
func (c *MetricDeleteCommand) Run(args []string) error {
	fs := flag.NewFlagSet("metric-delete", flag.ContinueOnError)
	name := fs.String("name", "", "Give a name of a metric you want to delete")
	id := fs.Int("id", -1, "Give an ID of a metric you want to delete")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" && *id == -1 {
		return errors.New("Name and ID could not be empty at the same time")
	}

	// here we create inmemory metric service. essentially an abstract controller
	ms := mem.NewMetricService()
	if *name != "" {
		err := ms.DeleteMetricByName(*name)
		return err
	} else {
		err := ms.DeleteMetricByID(*id)
		return err
	}
}
