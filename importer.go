package sap_segmentation

import (
	"log"
	"time"

	"github.com/tdx/sap_segmentation/model"
)

type Loader interface {
	LoadNext() ([]model.Segmentation, error)
}

type Storer interface {
	Store(segs []model.Segmentation) error
}

func ImportSegments(log *log.Logger, loader Loader, storer Storer) {

	rows := 0
	start := time.Now()

	for {
		segs, err := loader.LoadNext()
		if err != nil {
			log.Printf("importer: load segmentations failed: %v\n", err)
			return
		}
		if len(segs) == 0 {
			// all loaded
			log.Printf("importer: loaded %d records, time: %s\n", rows, time.Since(start))
			return
		}

		if err = storer.Store(segs); err != nil {
			log.Printf("importer: store segmentations failed: %v\n", err)
			return
		}

		rows += len(segs)
	}
}
