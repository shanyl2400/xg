package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)
var (
	HttpTimeout = time.Second * 5
	BaseURL = "http://localhost:8088/"
)

type JSONRequest struct {
	URL string
	Method string
	Query map[string]string
	JSONBody	interface{}

	Token string
}

func (jr JSONRequest) DoRequest(ctx context.Context) ([]byte, error){
	queryStrPairs := make([]string, len(jr.Query))
	i := 0
	for k, v := range jr.Query{
		queryStrPairs[i] = k + "=" + v
	}
	queryStr := ""
	if len(queryStrPairs) > 0 {
		queryStr = strings.Join(queryStrPairs, "&")
		queryStr = "?" + queryStr
	}

	reqBody := new(bytes.Buffer)
	if jr.JSONBody != nil {
		jsonData, err := json.Marshal(jr.JSONBody)
		if err != nil{
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}
	request, err := http.NewRequest(jr.Method, BaseURL + jr.URL + queryStr, reqBody)
	if err != nil{
		return nil, err
	}
	request.Header.Set("Content-type", "application/json")
	if jr.Token != "" {
		request.Header.Set("Authorization", jr.Token)
	}
	client := http.Client{
		Timeout: HttpTimeout,
	}

	resp, err := client.Do(request)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	return resBody, nil
}