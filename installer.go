package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
)

func connectSsh() (*ssh.Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir, %v", err)
	}

	key, err := os.ReadFile(homeDir + "/.ssh/id_rsa")
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: "pi",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", "192.168.1.107:22", config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to address: %v", err)
	}

	return conn, nil
}

func sendCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	var buff bytes.Buffer
	session.Stdout = &buff
	if err := session.Run(command); err != nil {
		return "", err
	}

	return buff.String(), nil
}

func build() {
	cmd := exec.Command("go", "build", "-o", "service", "-ldflags", "-s -w", "-installsuffix", "cgo", ".")
	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "GOOS=linux", "GOARCH=arm", "GOARM=7")
	if err := cmd.Run(); err != nil {
		log.Print(out.String())
		log.Fatal(err)
	}
}

func main2() {
	flag.Parse()
	build()
	log.Print(flag.Arg(0))
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := connectSsh()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client, err := scp.NewClientBySSH(conn)
	if err != nil {
		log.Fatal("Error creating new SSH session from existing connection", err)
	}

	fd, err := os.Open(wd + "/service")
	if err != nil {
		log.Fatal("Failed to read file", err)
	}

	if err := client.CopyFile(context.Background(), fd, flag.Arg(0), "0755"); err != nil {
		log.Fatal("Failed to copy file", err)
	}

	if _, err := sendCommand(conn, fmt.Sprintf("sudo mv ~/%s /apps/%s", flag.Arg(0), flag.Arg(0))); err != nil {
		log.Fatal("Failed to send command", err)
	}

	if _, err = sendCommand(conn, fmt.Sprintf("sudo systemctl restart %s.service", flag.Arg(1))); err != nil {
		log.Fatal("Failed to restart service", err)
	}
	if err := os.Remove(wd + "/service"); err != nil {
		log.Fatal("Failed to delete file", err)
	}
}
