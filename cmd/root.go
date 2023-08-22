/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/fr0stylo/imployer/pkg/build"
	"github.com/fr0stylo/imployer/pkg/config"
	"github.com/fr0stylo/imployer/pkg/deploy"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "imployer",
	Short: "This cli tool is used for application build and deployment",
	Long:  `This tool allows quick deployments into several machines via ssh that are running systemd services`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Executing build\n")

		cfg, err := config.Load(cfgFile)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}

		if cfg.Build != nil {
			if err := build.Execute(cfg.Build); err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}

		if cfg.Deploy != nil {
			fmt.Print("Deploying\n")
			if err := deploy.Execute(cfg.Deploy); err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
		}

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	rootCmd.Flags().StringVarP(&cfgFile, "config", "f", path.Join(dir, "install.yaml"), "install config file path")
}
