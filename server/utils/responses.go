package utils

import (
	"encoding/json"
	"fmt"

	"github.com/redsift/go-sandbox-rpc"
	rpc "github.com/redsift/go-sandbox-rpc/rpc"
)

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
	return sandboxrpc.NewComputeResponse("rpc_resp", k, v, 0, 0)
}

func ErrorResponse(key string, msg string, err error) sandboxrpc.ComputeResponse {
	et := err.Error()
	if len(msg) > 0 {
		et = fmt.Sprintf("%s: %s\n", msg, err.Error())
	}
	return NewComputeResponse(key, 500, []byte(et), false)
}

func ExportStats(statsmap map[string]interface{}) sandboxrpc.ComputeResponse {
	idx_stats, _ := json.Marshal(statsmap)
	return sandboxrpc.NewComputeResponse("stats", "index_stats", idx_stats, 0, 0)
}
