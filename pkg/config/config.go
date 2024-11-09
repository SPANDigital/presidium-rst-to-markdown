package config

import (
	"flag"
	"os"
)

const (
	DirPermission  = 0755
	FilePermission = 0644
)

type Config struct {
	InputDir    string
	OutputDir   string
	PandocPath  string
	Force       bool
	Verbose     bool
	MaxParallel int
	Depth       int // Maximum heading depth to split sections
}

// ParseArgs parses command-line arguments and returns a Config struct.
func ParseArgs() Config {
	var config Config
	flag.StringVar(&config.InputDir, "input", "", "Input directory")
	flag.StringVar(&config.OutputDir, "output", "", "Output directory")
	flag.StringVar(&config.PandocPath, "pandoc-path", "pandoc", "Path to the Pandoc executable")
	flag.BoolVar(&config.Force, "force", false, "Force overwrite of output directory")
	flag.BoolVar(&config.Verbose, "v", false, "Enable verbose logging")
	flag.IntVar(&config.MaxParallel, "parallel", 4, "Maximum number of parallel processes")
	flag.IntVar(&config.Depth, "depth", 2, "Heading depth level to split sections")

	flag.Parse()

	if config.InputDir == "" || config.OutputDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	return config
}
