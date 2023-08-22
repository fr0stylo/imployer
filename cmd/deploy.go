package cmd

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"

	"github.com/fr0stylo/imployer/pkg/deploy"
	"github.com/fr0stylo/imployer/pkg/utils"
)

var profiles = []string{}
var deleteArtifact bool = false

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy to server",
	Long:  `Deploy to server`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flag("serviceName").NoOptDefVal = args[0]
		cmd.Flag("executable").NoOptDefVal = args[0]

		fmt.Print("Deploying\n")
		sshOpts := lo.Map(profiles, func(v string, i int) *deploy.SshOpts {
			return &deploy.SshOpts{SshProfile: &v}
		})

		if len(profiles) == 0 {
			sshOpts = []*deploy.SshOpts{{
				Host:       utils.String(cmd.Flag("host").Value.String()),
				Port:       utils.String(cmd.Flag("port").Value.String()),
				User:       utils.String(cmd.Flag("user").Value.String()),
				PrivateKey: utils.String(cmd.Flag("identityFile").Value.String()),
			}}
		}

		if err := deploy.Execute(&deploy.DeployOpts{
			DeployablePath:      args[0],
			ExecutableName:      utils.String(cmd.Flag("executable").Value.String()),
			RemoteExecutableDir: utils.String(cmd.Flag("remoteExecutableDir").Value.String()),
			ServiceName:         utils.String(cmd.Flag("serviceName").Value.String()),
			SshOpts:             sshOpts,
			DeleteArtifact:      deleteArtifact,
		}); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringP("executable", "e", "", "")
	deployCmd.Flags().StringP("remoteExecutableDir", "r", "/apps", "")
	deployCmd.Flags().StringP("serviceName", "s", "", "")
	deployCmd.Flags().BoolVarP(&deleteArtifact, "deleteArtifact", "d", false, "")

	deployCmd.Flags().StringArrayVarP(&profiles, "profile", "p", []string{}, "")

	deployCmd.Flags().StringP("host", "H", "", "")
	deployCmd.Flags().StringP("port", "P", "22", "")
	deployCmd.Flags().StringP("user", "u", "user", "")
	deployCmd.Flags().StringP("identityFile", "i", "~/.ssh/id_rsa", "")

	deployCmd.MarkFlagsMutuallyExclusive("profile", "host")
	deployCmd.MarkFlagsRequiredTogether("host", "port", "user", "identityFile")
}
