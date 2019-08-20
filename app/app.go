package app

import (
	"fmt"
	"github.com/spaiz/mongotrace/config"
	"github.com/spaiz/mongotrace/db"
	"github.com/spf13/cobra"
	"log"
	"time"
)

// Options
type Options struct {
	Level  int
	Info   bool
	Config string
	Raw    bool
	Color  bool
	Indent bool
	Debug bool
}

// Returns App instance
func NewApp(version string) *App {
	return &App{
		version: version,
		options: &Options{},
	}
}

// Defines App structure
type App struct {
	version string
	conf    *config.Configuration
	options *Options
}

// Loads configuration file
func (r *App) loadConfig() error {
	var err error
	r.conf, err = config.LoadConfig(r.options.Config)
	return err
}

// Returns New database instances
func (r *App) getDatabases() []*db.MongoDB {
	hostsNum := len(r.conf.Hosts)
	dbs := make([]*db.MongoDB, 0, hostsNum)
	timeout := 20 * time.Second

	for _, dbHost := range r.conf.Hosts {
		conn := db.NewMongoDB(dbHost, timeout)

		err := conn.Connect()
		if err != nil {
			log.Printf("Failed to connect to db; %v\n", err.Error())
			continue
		}

		err = conn.Ping()
		if err != nil {
			log.Printf("Failed to ping db; %v\n", err.Error())
			continue
		}

		dbs = append(dbs, conn)
	}

	return dbs
}

// Starts application execution
func (r *App) Run() error {
	var rootCmd = &cobra.Command{
		Use:   "mongotrace",
		Short: "Tool for printing mongodb query logs",
		Long:  `Tool for printing mongodb query logs from multiple servers simultaneously`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
	}

	rootCmd.PersistentFlags().StringVar(&r.options.Config, "config", "./confs/config.json", "path to config file")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("mongotrace v%s\n", r.version)
		},
	}

	var enableCmd = &cobra.Command{
		Use:   "enable",
		Short: "Enables logs on the servers",
		Long:  `Enables logs on the servers by setting logs level to be 2"`,
		Run: func(cmd *cobra.Command, args []string) {
			r.options.Level = 2
			dbs := r.getDatabases()
			command := NewSetCommand(r.conf, r.options, dbs)
			err := command.Execute()
			if err != nil {
				log.Fatalf("Failed to run %s command: %s\n", cmd.Use, err.Error())
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			r.loadConfig()
		},
	}

	var disableCmd = &cobra.Command{
		Use:   "disable",
		Short: "Disables logs on the servers",
		Long:  `Disables logs on the servers by setting logs level to be 0"`,
		Run: func(cmd *cobra.Command, args []string) {
			dbs := r.getDatabases()
			r.options.Level = 0
			command := NewSetCommand(r.conf, r.options, dbs)
			err := command.Execute()
			if err != nil {
				log.Fatalf("Failed to run %s command: %s\n", cmd.Use, err.Error())
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			r.loadConfig()
		},
	}

	var infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Prints current logs level status on the hosts",
		Run: func(cmd *cobra.Command, args []string) {
			dbs := r.getDatabases()
			command := NewInfoCommand(r.conf, dbs)
			err := command.Execute()
			if err != nil {
				log.Fatalf("Failed to run %s command: %s\n", cmd.Use, err.Error())
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			r.loadConfig()
		},
	}

	var tailCmd = &cobra.Command{
		Use:   "tail",
		Short: "Prints raw mongodb query logs",
		Long:  `Prints raw mongodb query logs. You can use additional flags to enable formatting, colors and etc.`,
		Run: func(cmd *cobra.Command, args []string) {
			dbs := r.getDatabases()
			formatter := NewDefaultOplogFormatter()
			command := NewLogReportCommand(r.conf, r.options, dbs, formatter)
			err := command.Execute()
			if err != nil {
				log.Fatalf("Failed to run %s command: %s\n", cmd.Use, err.Error())
			}
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			r.loadConfig()
		},
	}

	tailCmd.Flags().BoolVarP(&r.options.Raw, "raw", "r", false, "Forces to print raw data from oplog")
	tailCmd.Flags().BoolVarP(&r.options.Color, "color", "c", true, "Enables syntax highlight for the raw data")
	tailCmd.Flags().BoolVarP(&r.options.Indent, "indent", "i", true, "Enables JSON indentation for the raw data")
	tailCmd.Flags().BoolVarP(&r.options.Debug, "debug", "d", false, "Prints raw data and formatted version")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(enableCmd)
	rootCmd.AddCommand(disableCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(tailCmd)

	return rootCmd.Execute()
}