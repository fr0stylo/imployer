package deploy

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"

	"github.com/fr0stylo/imployer/pkg/utils"
)

type SshOpts struct {
	SshProfile *string `yaml:"profile,omitempty"`
	Host       *string `yaml:"host,omitempty"`
	Port       *string `yaml:"port,omitempty"`
	PrivateKey *string `yaml:"privateKey,omitempty"`
	User       *string `yaml:"user,omitempty"`
}

func (r *SshOpts) Build() *SshOpts {
	if r.SshProfile != nil {
		r.Host = utils.String(ssh_config.Get(*r.SshProfile, "Hostname"))
		r.Port = utils.String(ssh_config.Get(*r.SshProfile, "Port"))
		r.User = utils.String(ssh_config.Get(*r.SshProfile, "User"))

		if ssh_config.Get(*r.SshProfile, "IdentityFile") == "~/.ssh/identity" {
			r.PrivateKey = utils.String("~/.ssh/id_rsa")
		} else {
			r.PrivateKey = utils.String(ssh_config.Get(*r.SshProfile, "IdentityFile"))
		}
	}

	return r
}

func connectSsh(opts *SshOpts) (*ssh.Client, error) {
	if opts.SshProfile != nil {
		fmt.Printf("Connecting to %s profile\n", *opts.SshProfile)
	} else {
		fmt.Printf("Connecting to %s host\n", *opts.Host)
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir, %v", err)
	}
	key, err := os.ReadFile(path.Join(strings.Replace(*opts.PrivateKey, "~", homeDir, 1)))
	// key, err := os.ReadFile(homeDir + "/.ssh/id_rsa")
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: *opts.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", *opts.Host, *opts.Port), config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to address: %v", err)
	}

	return conn, nil
}

func sendCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var buff bytes.Buffer
	session.Stdout = &buff
	if err := session.Run(command); err != nil {
		return "", err
	}

	return buff.String(), nil
}
