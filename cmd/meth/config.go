package main

import (
	"fmt"
	"os"

	"github.com/demolemo/meth/config"
	"github.com/demolemo/meth/sqlite"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show where meth keeps its config and its data",
	Args:  cobra.NoArgs,
	// cobra runs only the nearest PersistentPreRunE, so this empty one keeps
	// the root's database wiring out of the way: asking where the data lives
	// must not be what creates it
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath, err := config.Path()
		if err != nil {
			return err
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Printf("config file  %s\n", cfgPath)
		if cfg.File == "" {
			fmt.Printf("             not created yet, using defaults. write it with: meth config init\n")
		}

		fmt.Printf("database     %s\n", cfg.DB)
		if dbFlag != "" {
			fmt.Printf("             from the --db flag\n")
		}
		fmt.Printf("             %s\n", describeDB(cfg.DB))
		return nil
	},
}

// describeDB summarises the database without ever creating it.
func describeDB(path string) string {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "not created yet, the first command that writes will create it"
	}
	if err != nil {
		return err.Error()
	}

	db, err := sqlite.Open(path)
	if err != nil {
		return err.Error()
	}
	defer db.Close()

	mets, err := sqlite.NewMetricService(db).GetMetrics()
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("%d metrics, %.1f KB", len(mets), float64(info.Size())/1024)
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Write a config file with the current settings",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		path, err := config.Write(cfg)
		if os.IsExist(err) {
			return fmt.Errorf("%s already exists, edit it instead", path)
		}
		if err != nil {
			return err
		}
		fmt.Printf("Wrote %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
}
