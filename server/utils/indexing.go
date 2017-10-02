package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/redsift/blevex/rocksdb"
)

const BatchSize = 1000
const LinesFromCSV = 1000

type MajesticDatum struct {
	GlobalRank     string
	TldRank        string
	Domain         string
	TLD            string
	RefSubNets     string
	RefIPs         string
	IDN_Domain     string
	IDN_TLD        string
	PrevGlobalRank string
	PrevTldRank    string
	PrevRefSubNets string
	PrevRefIPs     string
}

func OpenIndex(forSearch bool) (bleve.Index, error) {
	_ = rocksdb.Name
	indexPath := os.Getenv("_LARGE_STORAGE_rocksdb_store")

	var idx bleve.Index
	var err error
	start := time.Now()
	if forSearch {
		cfg := map[string]interface{}{
			"read_only": true,
		}
		idx, err = bleve.OpenUsing(indexPath, cfg)
		if err != nil {
			err = errors.New("Tried to search before indexing: " + err.Error())
		}
	} else {
		idx, err = bleve.Open(indexPath)
		if err != nil {
			// make a new one
			fmt.Println("Creating a new index!")

			stdStoredAndIndexed := bleve.NewTextFieldMapping()
			stdStoredAndIndexed.Store = true
			stdStoredAndIndexed.IncludeInAll = true
			stdStoredAndIndexed.IncludeTermVectors = false

			numIndexed := bleve.NewNumericFieldMapping()
			numIndexed.Store = true
			numIndexed.IncludeInAll = true

			lineMapping := bleve.NewDocumentMapping()
			lineMapping.AddFieldMappingsAt("GlobalRank", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("TldRank", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("Domain", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("TLD", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("RefSubNets", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("RefIPs", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("IDN_Domain", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("IDN_TLD", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("PrevGlobalRank", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("PrevTldRank", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("PrevRefSubNets", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("PrevRefIPs", stdStoredAndIndexed)

			mapping := bleve.NewIndexMapping()
			mapping.DefaultMapping = lineMapping
			mapping.DefaultAnalyzer = "standard"
			idx, err = bleve.New(indexPath, mapping)

		}
	}
	if err != nil {
		return nil, err
	}
	t := "write"
	if forSearch {
		t = "read"
	}
	fmt.Printf("Index opened in %s (%s)\n", time.Now().Sub(start), t)
	return idx, nil
}

func UpdateIndex(idx bleve.Index, lines [][]string) error {
	start := time.Now()

	var batch *bleve.Batch

	for i, s := range lines {
		if i == LinesFromCSV {
			break
		}
		if batch == nil {
			batch = idx.NewBatch()
		}

		mjl := MajesticDatum{s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11]}
		if len(mjl.GlobalRank) == 0 {
			continue
		}
		if err := batch.Index(mjl.GlobalRank, mjl); err != nil {
			return err
		}

		if batch.Size() == BatchSize {
			fmt.Println("committing batch!", i)
			if err := idx.Batch(batch); err != nil {
				return err
			}
			batch = nil
		}
	}

	if batch != nil {
		if err := idx.Batch(batch); err != nil {
			return err
		}
		batch = nil
	}

	fmt.Printf("Indexed %d lines in %0.3fs\n", LinesFromCSV, time.Now().Sub(start).Seconds())
	return nil
}
