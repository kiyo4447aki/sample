package main

import (
	"camera/config"
	"camera/streamer"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	sender, err := streamer.BuildPipeline(*cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := sender.Run(*cfg); err != nil {
		log.Fatal(err.Error())
	}

}