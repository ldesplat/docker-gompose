package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
)

func stringStartsInSlice(a string, list []string) bool {

	for _, b := range list {
		if strings.HasPrefix(b, a) {
			return true
		}
	}
	return false
}

func before(c *cli.Context) (Containers, *docker.Client, string, error) {

	var configFile = "docker-compose.yml"
	if c.GlobalString("file") != "FILE" {
		configFile = c.GlobalString("file")
	}

	var workingDir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	var projectName = path.Base(workingDir)

	if c.GlobalString("file") != "FILE" {
		var baseDir = path.Dir(configFile)
		if baseDir != "." && baseDir != "/" {
			projectName = strings.ToLower(path.Base(baseDir))
		}
	}
	if c.GlobalString("project-name") != "NAME" {
		projectName = c.GlobalString("project-name")
	}

	// Parse the yml file
	config, err := ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Now we connect to the docker daemon
	client, err := ConnectToDocker()
	if err != nil {
		log.Fatal(err)
	}

	return config, client, projectName, err
}

func main() {
	cli.AppHelpTemplate = `{{.Usage}}

Usage:
  docker-gompose [options] [COMMAND] [ARGS...]

Options:
  {{range .Flags}}{{.}}
  {{end}}

Commands:
  {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
  {{end}}
`
	app := cli.NewApp()
	app.Name = "docker-gompose"
	app.Usage = "Fast, isolated development environments using Docker."
	app.Version = "1.2.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "file, f",
			Value:  "FILE",
			Usage:  "Specify an alternate compose file (default: docker-compose.yml)",
			EnvVar: "COMPOSE_FILE",
		},
		cli.StringFlag{
			Name:   "project-name, p",
			Value:  "NAME",
			Usage:  "Specify an alternate project name (default: directory name)",
			EnvVar: "COMPOSE_PROJECT_NAME",
		},
		cli.StringFlag{
			Name:  "verbose",
			Usage: "Show more output",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "build",
			Usage: "Build or rebuild services",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "help",
			Usage: "Get help on a command",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "kill",
			Usage: "Kill containers",
			Action: func(c *cli.Context) {
				config, client, projectName, _ := before(c)
				CmdKill(config, client, projectName)
			},
		},
		{
			Name:  "logs",
			Usage: "View output from containers",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "port",
			Usage: "Print the public port for a port binding",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "ps",
			Usage: "List containers",
			Action: func(c *cli.Context) {
				config, client, projectName, _ := before(c)
				CmdPs(config, client, projectName)
			},
		},
		{
			Name:  "pull",
			Usage: "Pulls service images",
			Action: func(c *cli.Context) {
				config, client, _, _ := before(c)
				CmdPull(config, client)
			},
		},
		{
			Name:  "rm",
			Usage: "Remove stopped containers",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "run",
			Usage: "Run a one-off command",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "scale",
			Usage: "Set number of containers for a service",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "start",
			Usage: "Start services",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "stop",
			Usage: "Stop services",
			Action: func(c *cli.Context) {
				config, client, projectName, _ := before(c)
				CmdStop(config, client, projectName)
			},
		},
		{
			Name:  "restart",
			Usage: "Restart services",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
		{
			Name:  "up",
			Usage: "Create and start containers",
			Action: func(c *cli.Context) {
				fmt.Println("Not yet implemented!")
			},
		},
	}

	app.Run(os.Args)
}
