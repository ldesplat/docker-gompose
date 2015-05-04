package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/tabwriter"

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

// CmdPs defines the ps command
func CmdPs(config Containers, client *docker.Client, projectName string) {

	conts, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatalf("%v", err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 1, 3, ' ', 0)
	fmt.Fprintln(w, "Name\tCommand\tState\tPorts")
	fmt.Fprintln(w, "-----------------\t-----------------\t-----------------\t-----------------")
	for _, cont := range conts {
		if stringStartsInSlice("/"+projectName+"_", cont.Names) {
			//fmt.Printf("%-.18s   %-.18s   %-.18s   ", cont.Names[0][1:len(cont.Names[0])], cont.Command, cont.Status)
			fmt.Fprintf(w, "%s\t%s\t%-.18s\t", cont.Names[0][1:len(cont.Names[0])], cont.Command, cont.Status)
			for i, port := range cont.Ports {
				if i > 0 {
					fmt.Fprintf(w, "\n\t\t\t")
				}
				fmt.Fprintf(w, "%s:%v->%v/%s", port.IP, port.PublicPort, port.PrivatePort, port.Type)
			}

			fmt.Fprintln(w)
		}
	}
	w.Flush()
}

// CmdPull defines the pull command
func CmdPull(config Containers, client *docker.Client) {

	for name, cont := range config {
		if cont.Image != "" {

			fmt.Printf("Pulling %s (%s)...\n", name, cont.Image)
			err := client.PullImage(docker.PullImageOptions{Repository: cont.Image}, docker.AuthConfiguration{})
			if err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
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
			Name:  "file, f",
			Value: "FILE",
			Usage: "Specify an alternate compose file (default: docker-compose.yml)",
		},
		cli.StringFlag{
			Name:  "project-name, p",
			Value: "NAME",
			Usage: "Specify an alternate project name (default: directory name)",
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
				fmt.Println("Build")
			},
		},
		{
			Name:  "help",
			Usage: "Get help on a command",
			Action: func(c *cli.Context) {
				fmt.Println("Help")
			},
		},
		{
			Name:  "kill",
			Usage: "Kill containers",
			Action: func(c *cli.Context) {
				fmt.Println("Kill")
			},
		},
		{
			Name:  "logs",
			Usage: "View output from containers",
			Action: func(c *cli.Context) {
				fmt.Println("Logs")
			},
		},
		{
			Name:  "port",
			Usage: "Print the public port for a port binding",
			Action: func(c *cli.Context) {
				fmt.Println("Port")
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
				fmt.Println("rm")
			},
		},
		{
			Name:  "run",
			Usage: "Run a one-off command",
			Action: func(c *cli.Context) {
				fmt.Println("run")
			},
		},
		{
			Name:  "scale",
			Usage: "Set number of containers for a service",
			Action: func(c *cli.Context) {
				fmt.Println("scale")
			},
		},
		{
			Name:  "start",
			Usage: "Start services",
			Action: func(c *cli.Context) {
				fmt.Println("start")
			},
		},
		{
			Name:  "stop",
			Usage: "Stop services",
			Action: func(c *cli.Context) {
				fmt.Println("stop")
			},
		},
		{
			Name:  "restart",
			Usage: "Restart services",
			Action: func(c *cli.Context) {
				fmt.Println("Build")
			},
		},
		{
			Name:  "up",
			Usage: "Create and start containers",
			Action: func(c *cli.Context) {
				fmt.Println("up")
			},
		},
	}

	app.Run(os.Args)
}
