package deploy

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/bramvdbogaerde/go-scp"
)

type DeployOpts struct {
	DeployablePath      string     `yaml:"input"`
	ExecutableName      *string    `yaml:"execName,omitempty"`
	RemoteExecutableDir *string    `yaml:"remoteDir,omitempty"`
	ServiceName         *string    `yaml:"service,omitempty"`
	DeleteArtifact      bool       `yaml:"delete,omitempty"`
	SshOpts             []*SshOpts `yaml:"ssh,omitempty"`
}

func Execute(opts *DeployOpts) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, v := range opts.SshOpts {
		conn, err := connectSsh(v.Build())
		if err != nil {
			return err
		}
		defer conn.Close()

		client, err := scp.NewClientBySSH(conn)
		if err != nil {
			return fmt.Errorf("error creating new SSH session from existing connection, %v", err)
		}

		fd, err := os.Open(path.Join(wd, opts.DeployablePath))
		if err != nil {
			return fmt.Errorf("failed to read file, %s", err)
		}

		if err := client.CopyFile(context.Background(), fd, *opts.ExecutableName, "0755"); err != nil {
			return fmt.Errorf("failed to copy file %v", err)
		}

		if err := fd.Close(); err != nil {
			return fmt.Errorf("failed to close file %v", err)
		}

		if _, err := sendCommand(conn,
			fmt.Sprintf("sudo mv ~/%s %s",
				*opts.ExecutableName,
				path.Join(*opts.RemoteExecutableDir, *opts.ExecutableName))); err != nil {
			return fmt.Errorf("failed to send command %v", err)
		}

		if _, err = sendCommand(conn, fmt.Sprintf("sudo systemctl restart %s.service", *opts.ServiceName)); err != nil {
			return fmt.Errorf("failed to restart service, %s", err)
		}
	}

	if opts.DeleteArtifact {
		if err := os.Remove(path.Join(wd, opts.DeployablePath)); err != nil {
			return fmt.Errorf("failed to delete file %v", err)
		}
	}

	return nil
}
