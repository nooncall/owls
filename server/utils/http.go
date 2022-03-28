package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	PostMethod = "POST"
	GetMethod  = "GET"
)

func DoHttpReq(method, url, jsonParam string, header http.Header) ([]byte, error) {
	cli := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(jsonParam))
	if err != nil {
		return nil, fmt.Errorf("new request err:%s", err.Error())
	}

	req.Header = header
	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do  request err: %s", err.Error())
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
