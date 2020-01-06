package search

import (
	"encoding/json"
	"fmt"
	"net/url"
	"server/utils"
	"strings"

	"github.com/redsift/bleve"
	"github.com/redsift/go-sandbox-rpc"
	rpc "github.com/redsift/go-sandbox-rpc/rpc"
)

func Compute(req sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error) {
	inData := req.In.Data
	if len(inData) != 1 {
		return nil, fmt.Errorf("empty input")
	}

	var resp []sandboxrpc.ComputeResponse
	v := inData[0]
	idx, err := utils.OpenIndex(true)
	if err != nil {
		resp = append(resp, utils.ErrorResponse(v.Data.Key, "error creating index", err))
		return resp, nil
	}
	defer idx.Close()

	var rpcReq rpc.Request
	err = json.Unmarshal(v.Data.Value, &rpcReq)
	if err != nil {
		resp = append(resp, utils.ErrorResponse(v.Data.Key, "Unmarshal the rpc request failed", err))
		return resp, nil
	}

	u, err := url.ParseRequestURI(rpcReq.RequestURI)
	if err != nil {
		resp = append(resp, utils.ErrorResponse(v.Data.Key, "Something went wrong with parsing RequestURI", err))
		return resp, nil
	}

	htmlQuery := u.Query()
	q := strings.TrimSpace(htmlQuery.Get("q"))
	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
	searchRequest.Fields = []string{"*"}
	searchResult, err := idx.Search(searchRequest)
	if err != nil {
		resp = append(resp, utils.ErrorResponse(v.Data.Key, "", err))
		return resp, nil
	}

	b, _ := json.Marshal(searchResult)
	resp = append(resp, utils.NewComputeResponse(v.Data.Key, 200, b, true))
	resp = append(resp, utils.ExportStats(idx))

	return resp, nil
}
