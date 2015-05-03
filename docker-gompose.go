package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/fsouza/go-dockerclient"
)

var data = `
web:
  ports:
    - "8000:8000"
elasticsearch:
  image: elasticsearch:latest
  ports:
    - "9200:9200"
    - "9300:9300"
neo4j:
  image: tpires/neo4j
  ports:
    - "7474:7474"
    - "1337:1337"
  external_links:
    - bla
    - bla2
  env_file: yellow
redis:
  image: redis:latest
  PORTS:
    - "6379:6379"
  env_file:
    - green
    - red
    - blue
nginx:
  image: nginx:latest
`

func before(c *cli.Context) error {

	// Parse the yml file
	config, err := ParseConfig("./docker-compose.yml")
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("--- t:\n%v\n\n", config)

	for index, element := range config {
		fmt.Printf("------ %v ------\n", index)
		fmt.Printf("Image: %v\n", element.Image)
		fmt.Printf("Ports: %v\n", element.Ports)
		fmt.Printf("Volume: %v\n", element.Volumes)
		fmt.Printf("Links: %v\n", element.Links)
	}

	// Now we connect to the docker daemon
	client, err := ConnectToDocker()
	if err != nil {
		log.Fatal(err)
		return err
	}

	images, err := client.ListImages(docker.ListImagesOptions{All: true})

	if err != nil {
		log.Fatal(err)
	}

	for _, c := range images {
		log.Println(c.ID, c.Created)
	}

	return err
}

func main() {
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
			Name:   "build",
			Usage:  "Builds or rebuilds services",
			Before: before,
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
				fmt.Println("ps")
			},
		},
		{
			Name:  "pull",
			Usage: "Pulls service images",
			Action: func(c *cli.Context) {
				fmt.Println("Pull")
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
