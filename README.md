# Imployer

Imployer (Employer) - CLI tool for building and deploying applications to remote machines via SSH. This tool is capable of deploying single binary to several remote machines and restarting systemd service afterwards.

## Install
### Remote
In order to install as cli run:
`go install github.com/fr0stylo/imployer@latest`

### Building from source

Clone repository and run `go install .`

## Usage

```
Usage:
  imployer [flags]
  imployer [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  deploy      Deploy to server
  help        Help about any command

Flags:
  -f, --config string   install config file path (default "./install.yaml")


Deploy to server

Usage:
  imployer deploy [flags]

Flags:
  -d, --deleteArtifact               
  -e, --executable string            
  -h, --help                         help for deploy
  -H, --host string                  
  -i, --identityFile string           (default "~/.ssh/id_rsa")
  -P, --port string                   (default "22")
  -p, --profile stringArray          
  -r, --remoteExecutableDir string    (default "/apps")
  -s, --serviceName string           
  -u, --user string                   (default "user")

```

install.yaml configuration file example:

```
build:
  env:
    - "GOOS=linux"
    - "GOARCH=arm"
    - "GOARM=7"
  flags:
    - "-ldflags"
    - "-s -w"
    - "-installsuffix"
    - "cgo"
  output: service
deploy:
  input: service
  execName: application
  remoteDir: /apps
  service: context-application
  delete: true
  ssh:
    - profile: pi

```

## Contribution

All contributions are welcome as pull request
