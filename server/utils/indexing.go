package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/redsift/blevex/rocksdb"
)

const BatchSize = 1000
const LinesFromCSV = 1000

type MajesticDatum struct {
	GlobalRank     int
	TldRank        int
	Domain         string
	TLD            string
	RefSubNets     int
	RefIPs         int
	IDN_Domain     string
	IDN_TLD        string
	PrevGlobalRank int
	PrevTldRank    int
	PrevRefSubNets int
	PrevRefIPs     int
}

func newMajecticDatum(line []string) MajesticDatum {
	ls := make([]int, 12)
	for i, value := range line {
		t, _ := strconv.Atoi(value)
		ls[i] = t
	}
	return MajesticDatum{ls[0], ls[1], line[2], line[3], ls[4], ls[5], line[6], line[7], ls[8], ls[9], ls[10], ls[11]}
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
			lineMapping.AddFieldMappingsAt("GlobalRank", numIndexed)
			lineMapping.AddFieldMappingsAt("TldRank", numIndexed)
			lineMapping.AddFieldMappingsAt("Domain", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("TLD", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("RefSubNets", numIndexed)
			lineMapping.AddFieldMappingsAt("RefIPs", numIndexed)
			lineMapping.AddFieldMappingsAt("IDN_Domain", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("IDN_TLD", stdStoredAndIndexed)
			lineMapping.AddFieldMappingsAt("PrevGlobalRank", numIndexed)
			lineMapping.AddFieldMappingsAt("PrevTldRank", numIndexed)
			lineMapping.AddFieldMappingsAt("PrevRefSubNets", numIndexed)
			lineMapping.AddFieldMappingsAt("PrevRefIPs", numIndexed)

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

		mjl := newMajecticDatum(s)
		if err := batch.Index(strconv.Itoa(mjl.GlobalRank), mjl); err != nil {
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
