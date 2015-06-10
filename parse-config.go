package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// SingleOrMulti refers to a parameter that can include either a string or an array of strings
type SingleOrMulti struct {
	Values []string
}

func (sm *SingleOrMulti) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			return err
		}
		sm.Values = make([]string, 1)
		sm.Values[0] = single
	} else {
		sm.Values = multi
	}
	return nil
}

// ExtendConfig refers to a file and service that we must parse as well
type ExtendConfig struct {
	File    string
	Service string
}

// Container is the full representation of a container in the configuration
type Container struct {
	Image         string
	Build         string
	Command       string
	Links         []string
	ExternalLinks []string `yaml:"external_links"`
	Ports         []string
	Expose        []string
	Volumes       []string
	VolumesFrom   []string `yaml:"volumes_from"`
	Environment   []string
	EnvFile       SingleOrMulti `yaml:"env_file"`
	Extends       ExtendConfig
	Net           string
	PID           string
	DNS           SingleOrMulti
	CapAdd        []string      `yaml:"cap_add"`
	CapDrop       []string      `yaml:"cap_drop"`
	DNSSearch     SingleOrMulti `yaml:"dns_search"`
	WorkingDir    string        `yaml:"working_dir"`
	Entrypoint    string
	User          string
	Hostname      string
	Domainname    string
	MemLimit      string `yaml:"mem_limit"`
	Privileged    string
	Restart       string
	StdinOpen     string `yaml:"stdin_open"`
	TTY           string
	CPUShares     string `yaml:"cpu_shares"`
}

// Containers includes all containers specified in configuration
type Containers map[string]Container

// ParseConfig parses a specified *.yml file
func ParseConfig(file string) (Containers, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := Containers{}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	return config, err
}
