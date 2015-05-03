package main

import (
	"fmt"
	"os"

	"github.com/fsouza/go-dockerclient"
)

// ConnectToDocker retrieves a connection to the docker daemon
func ConnectToDocker() (*docker.Client, error) {
	path := os.Getenv("DOCKER_CERT_PATH")
	client, err := docker.NewTLSClient(os.Getenv("DOCKER_HOST"), fmt.Sprintf("%s/cert.pem", path), fmt.Sprintf("%s/key.pem", path), fmt.Sprintf("%s/ca.pem", path))

	if err != nil {
		return nil, err
	}

	return client, err
}
