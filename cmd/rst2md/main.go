package main

import (
	"io"
	"log"

	"github.com/spandigital/presidium-rst-to-markdown/pkg/config"
	"github.com/spandigital/presidium-rst-to-markdown/pkg/processor"
)

func main() {
	cfg := config.ParseArgs()

	if cfg.Verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(0)
		log.SetOutput(io.Discard)
	}

	if err := processor.Run(cfg); err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Println("Conversion completed successfully.")
}
