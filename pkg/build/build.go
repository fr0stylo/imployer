package build

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type BuildOpts struct {
	EnvironmentVariables []string `yaml:"env"`
	BuildFlags           []string `yaml:"flags"`
	OutputPath           string   `yaml:"output"`
}

func Execute(opts *BuildOpts) error {
	cmds := []string{"build"}
	cmds = append(cmds, opts.BuildFlags...)
	cmds = append(cmds, "-o", opts.OutputPath)
	cmds = append(cmds, ".")

	cmd := exec.Command("go", cmds...)

	var out bytes.Buffer
	cmd.Stderr = &out
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, opts.EnvironmentVariables...)

	if err := cmd.Run(); err != nil {
		fmt.Fprint(os.Stdout, out.String())
		return err
	}

	return nil
}
