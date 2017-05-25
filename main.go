package main

import (
	"flag"
	"log"

	v "github.com/appscode/go/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "kubed",
		Short: `Run Kubernetes Daemon for AppsCode`,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})
		},
	}
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	rootCmd.AddCommand(NewCmdRun())
	rootCmd.AddCommand(v.NewCmdVersion())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
