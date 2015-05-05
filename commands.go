package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fsouza/go-dockerclient"
)

var signalMap = map[string]docker.Signal{
	"SIGABRT":   docker.SIGABRT,
	"SIGALRM":   docker.SIGALRM,
	"SIGBUS":    docker.SIGBUS,
	"SIGCHLD":   docker.SIGCHLD,
	"SIGCLD":    docker.SIGCLD,
	"SIGCONT":   docker.SIGCONT,
	"SIGFPE":    docker.SIGFPE,
	"SIGHUP":    docker.SIGHUP,
	"SIGILL":    docker.SIGILL,
	"SIGINT":    docker.SIGINT,
	"SIGIO":     docker.SIGIO,
	"SIGIOT":    docker.SIGIOT,
	"SIGKILL":   docker.SIGKILL,
	"SIGPIPE":   docker.SIGPIPE,
	"SIGPOLL":   docker.SIGPOLL,
	"SIGPROF":   docker.SIGPROF,
	"SIGPWR":    docker.SIGPWR,
	"SIGQUIT":   docker.SIGQUIT,
	"SIGSEGV":   docker.SIGSEGV,
	"SIGSTKFLT": docker.SIGSTKFLT,
	"SIGSTOP":   docker.SIGSTOP,
	"SIGSYS":    docker.SIGSYS,
	"SIGTERM":   docker.SIGTERM,
	"SIGTRAP":   docker.SIGTRAP,
	"SIGTSTP":   docker.SIGTSTP,
	"SIGTTIN":   docker.SIGTTIN,
	"SIGTTOU":   docker.SIGTTOU,
	"SIGUNUSED": docker.SIGUNUSED,
	"SIGURG":    docker.SIGURG,
	"SIGUSR1":   docker.SIGUSR1,
	"SIGUSR2":   docker.SIGUSR2,
	"SIGVTALRM": docker.SIGVTALRM,
	"SIGWINCH":  docker.SIGWINCH,
	"SIGXCPU":   docker.SIGXCPU,
	"SIGXFSZ":   docker.SIGXFSZ,
}

// CmdPs defines the ps command
func CmdPs(config Containers, client *docker.Client, projectName string, onlyIds bool) {

	conts, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		log.Fatalf("%v", err)
	}

	w := new(tabwriter.Writer)
	if !onlyIds {
		w.Init(os.Stdout, 0, 1, 3, ' ', 0)
		fmt.Fprintln(w, "Name\tCommand\tState\tPorts")
		fmt.Fprintln(w, "-----------------\t-----------------\t-----------------\t-----------------")
	}
	for _, cont := range conts {
		if stringStartsInSlice("/"+projectName+"_", cont.Names) {
			if onlyIds {
				fmt.Println(cont.ID)
			} else {
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
	}
	if !onlyIds {
		w.Flush()
	}
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

// CmdStop defines the stop command
func CmdStop(config Containers, client *docker.Client, projectName string) {

	conts, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, container := range conts {
		if stringStartsInSlice("/"+projectName+"_", container.Names) {
			fmt.Printf("Stopping %s...\n", container.Names[0])
			client.StopContainer(container.ID, 30)
		}
	}
}

// CmdKill defines the kill command
func CmdKill(config Containers, client *docker.Client, projectName string, signal string) {

	var signalCode = signalMap["SIGKILL"]
	if signal != "SIGNAL" {
		signalCode = signalMap[strings.ToUpper(signal)]
	}
	conts, err := client.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, container := range conts {
		if stringStartsInSlice("/"+projectName+"_", container.Names) {
			fmt.Printf("Killing %s...\n", container.Names[0])
			client.KillContainer(docker.KillContainerOptions{ID: container.ID, Signal: signalCode})
		}
	}
}
