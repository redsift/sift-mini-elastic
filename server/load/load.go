package load

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"sift/utils"
	"time"

	"github.com/redsift/go-sandbox-rpc"
)

const MAJESTIC_CSV_URL = "http://downloads.majestic.com/majestic_million.csv"

func fetchMajesticCSV() ([][]string, error) {
	start := time.Now()
	resp, err := http.Get(MAJESTIC_CSV_URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Fetched csv in %0.3fs\n", time.Now().Sub(start).Seconds())
	return records, nil
}

func Compute(req sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error) {
	idx, err := utils.OpenIndex(false)
	if err != nil {
		return nil, errors.New("Something went wrong while creating the index: " + err.Error())
	}
	defer idx.Close()

	datums, err := fetchMajesticCSV()
	if err != nil {
		return nil, err
	}

	fmt.Println("updating index started!")
	err = utils.UpdateIndex(idx, datums)
	if err != nil {
		return nil, err
	}

	return []sandboxrpc.ComputeResponse{utils.ExportStats(idx.StatsMap())}, nil
}
