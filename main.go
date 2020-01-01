package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kardianos/service"
)

const defaultPort = 8081

var logger service.Logger

type program struct {
	port int
}

func (p program) Start(s service.Service) error {
	go startHttp(p.port)
	return nil
}

func (p program) Stop(s service.Service) error {
	stopHttp()
	return nil
}

func main() {
	conf := &service.Config{
		Name:        "HnTopRss",
		DisplayName: "HN Top Links RSS",
		Description: "RSS feed for Hacker News Top Links",
	}
	prg := &program{defaultPort}

	s, err := service.New(prg, conf)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		var err error
		verb := os.Args[1]
		switch verb {
		case "install":
			err = s.Install()
			if err != nil {
				fmt.Printf("Failed to install: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" installed.\n", conf.DisplayName)
		case "uninstall":
			err = s.Uninstall()
			if err != nil {
				fmt.Printf("Failed to uninstall: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", conf.DisplayName)
		case "run":
			err = s.Run()
			if err != nil {
				fmt.Printf("Failed to run: %s\n", err)
				return
			}
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", conf.DisplayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", conf.DisplayName)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func startHttp(port int) {
	logger.Infof("Listening on port %d", port)
	http.HandleFunc("/", hnRssHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
func stopHttp() {
	logger.Info("Shutting down")
}
