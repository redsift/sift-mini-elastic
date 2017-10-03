package search

import (
	"encoding/json"
	"errors"
	"net/url"
	"sift/utils"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/redsift/go-sandbox-rpc"
	rpc "github.com/redsift/go-sandbox-rpc/rpc"
)

func Compute(req sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error) {
	inData := req.In.Data
	idx, errI := utils.OpenIndex(true)
	if errI != nil {
		errI = errors.New("Something went wrong while creating the index: " + errI.Error())
	}
	defer func() {
		if errI == nil {
			idx.Close()
		}
	}()

	var resp []sandboxrpc.ComputeResponse
	for _, v := range inData {
		if errI != nil {
			resp = append(resp, utils.ErrorResponse(v.Data.Key, "", errI))
			break
		}

		var rpcReq rpc.Request
		err := json.Unmarshal(v.Data.Value, &rpcReq)
		if err != nil {
			resp = append(resp, utils.ErrorResponse(v.Data.Key, "Unmarshal the rpc request failed", err))
			break
		}

		u, err := url.ParseRequestURI(rpcReq.RequestURI)
		if err != nil {
			resp = append(resp, utils.ErrorResponse(v.Data.Key, "Something went wrong with parsing RequestURI", err))
			break
		}

		htmlQuery := u.Query()
		q := strings.TrimSpace(htmlQuery.Get("q"))
		searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
		searchRequest.Fields = []string{"*"}
		searchResult, err := idx.Search(searchRequest)
		if err != nil {
			resp = append(resp, utils.ErrorResponse(v.Data.Key, "", err))
			break
		}

		b, _ := json.Marshal(searchResult)
		resp = append(resp, utils.NewComputeResponse(v.Data.Key, 200, b, true))
	}

	if idx != nil {
		resp = append(resp, utils.ExportStats(idx.StatsMap()))
	}
	return resp, nil
}
