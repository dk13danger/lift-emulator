package main

import (
	"flag"
	"log"

	"github.com/dk13danger/lift-emulator/server"
	"github.com/dk13danger/lift-emulator/service"
)

func main() {
	count := flag.Uint("c", 0, "floors count")
	height := flag.Uint("h", 0, "floor height")
	speed := flag.Uint("s", 0, "lift speed")
	delay := flag.Duration("d", 0, "doors delay")

	flag.Parse()

	lift, err := service.NewLift(*count, *height, *speed, *delay)
	if err != nil {
		log.Fatalf("Error while creating lift: %v", err)
	}

	registry := service.NewRegistry()
	srv := service.New(lift, registry)
	eventsCh := srv.Run()

	log.Fatalf("%v", server.Run(eventsCh))
}
