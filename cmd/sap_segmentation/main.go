package main

import (
	"time"

	ss "github.com/tdx/sap_segmentation"
)

func main() {
	cfg, err := ss.LoadConfig()
	if err != nil {
		panic(err)
	}

	logMaxAge := time.Hour * time.Duration(24*cfg.Log.CleanupMaxAgeDays)
	logger, err := ss.NewLogger(logMaxAge)
	if err != nil {
		panic(err)
	}

	loader, err := ss.NewHTTPLoader(logger, &cfg.LoaderArgs)
	if err != nil {
		logger.Fatalf("failed to create importer: %v\n", err)
	}

	storer, err := ss.NewPGStorer(&cfg.StorerArgs)
	if err != nil {
		logger.Fatalf("failed to create storer: %v\n", err)
	}

	ss.ImportSegments(logger, loader, storer)
}
