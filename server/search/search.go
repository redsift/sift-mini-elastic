package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sift/utils"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/redsift/go-sandbox-rpc"
	rpc "github.com/redsift/go-sandbox-rpc/rpc"
)

func Compute(req sandboxrpc.ComputeRequest) ([]sandboxrpc.ComputeResponse, error) {
	inData := req.In.Data
	idx, err := utils.OpenIndex(true)
	if err != nil {
		err = errors.New("Something went wrong while creating the index: " + err.Error())
	}
	defer func() {
		idx.Close()
	}()

	var resp []sandboxrpc.ComputeResponse
	searchFields := []string{"*"}
	for _, v := range inData {
		if err != nil {
			resp = append(resp, ErrorResponse(v.Data.Key, "", err))
			break
		}
		var rpcReq rpc.Request
		err := json.Unmarshal(v.Data.Value, &rpcReq)
		if err != nil {
			resp = append(resp, ErrorResponse(v.Data.Key, "Unmarshal the rpc request failed", err))
			break
		}

		u, err := url.ParseRequestURI(rpcReq.RequestURI)
		if err != nil {
			resp = append(resp, ErrorResponse(v.Data.Key, "Something went wrong with parsing RequestURI", err))
			break
		}

		htmlQuery := u.Query()

		sf := htmlQuery.Get("sf")
		if len(sf) > 0 {
			searchFields = strings.Split(sf, ",")
		}

		q := strings.TrimSpace(htmlQuery.Get("q"))
		searchRequest := bleve.NewSearchRequest(bleve.NewMatchQuery(q))
		searchRequest.Fields = searchFields
		searchResult, err := idx.Search(searchRequest)
		if err != nil {
			resp = append(resp, ErrorResponse(v.Data.Key, "", err))
			break
		}

		b, _ := json.Marshal(searchResult)
		r := NewComputeResponse(v.Data.Key, 200, b, true)
		resp = append(resp, r)
	}

	resp = append(resp, utils.ExportStats(idx))
	return resp, nil
}

func NewComputeResponse(k string, status int, b []byte, isJSON bool) sandboxrpc.ComputeResponse {
	r := rpc.Response{
		StatusCode: status,
		Header:     map[string][]string{},
		Body:       b,
	}
	if isJSON {
		r.Header.Add("Content-Type", "application/json")
	}

	v, _ := json.Marshal(r)
	return sandboxrpc.NewComputeResponse("rpc_rep", k, v, 0, 0)
}

func ErrorResponse(key string, msg string, err error) sandboxrpc.ComputeResponse {
	et := err.Error()
	if len(msg) > 0 {
		et = fmt.Sprintf("%s: %s\n", msg, err.Error())
	}
	return NewComputeResponse(key, 500, []byte(et), false)
}
