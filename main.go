package main

import (
	"bitbucket.org/kardianos/service"
	"fmt"
	"net/http"
	"os"
)

const port = 8081

var log service.Logger

func main() {
	var name = "HnTopRss"
	var displayName = "HN Top Links RSS"
	var desc = "RSS feed for Hacker News Top Links"

	var s, err = service.NewService(name, displayName, desc)
	log = s

	if err != nil {
		fmt.Printf("%s unable to start: %s", displayName, err)
		return
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
			fmt.Printf("Service \"%s\" installed.\n", displayName)
		case "remove":
			err = s.Remove()
			if err != nil {
				fmt.Printf("Failed to remove: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" removed.\n", displayName)
		case "run":
			startHttp()
		case "start":
			err = s.Start()
			if err != nil {
				fmt.Printf("Failed to start: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" started.\n", displayName)
		case "stop":
			err = s.Stop()
			if err != nil {
				fmt.Printf("Failed to stop: %s\n", err)
				return
			}
			fmt.Printf("Service \"%s\" stopped.\n", displayName)
		}
		return
	}
	err = s.Run(func() error {
		// start
		go startHttp()
		return nil
	}, func() error {
		// stop
		stopHttp()
		return nil
	})
	if err != nil {
		s.Error(err.Error())
	}
}

func startHttp() {
	log.Info("Listening on port %d", port)
	http.HandleFunc("/", hnRssHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
func stopHttp() {
	log.Info("Shutting down")
}
