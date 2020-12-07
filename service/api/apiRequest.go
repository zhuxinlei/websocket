package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

//鉴权返回结果

type Result struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}


func DoRequest(ctx context.Context, method, url string, headers map[string]string, params interface{}) (*Result, error) {

	jsb, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(%+v) failed: %s", params, err.Error())
	}
	if params == nil {
		jsb = []byte{}
	}
	log.Debugf("[SECRET] %s %s -> \n%s\n", method, url, string(jsb))

	req, err := http.NewRequest(method, url, bytes.NewReader(jsb))

	if err != nil {
		return nil, fmt.Errorf("new requrest (%s, %s, %s) failed: %s", method, url, string(jsb), err.Error())
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	client.Timeout = time.Minute

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read from resp body failed: %s", err.Error())
	}

	result := Result{}
	err = json.Unmarshal(bs, &result)
	if err != nil {
		fmt.Println(bs)
		return nil, fmt.Errorf("unmarshal %s failed: %s", string(bs), err.Error())
	}

	if result.Code != 200 {
		return &result, fmt.Errorf("method %s url %s header %+v params %+v error code [%d] %s",
			method, url, headers, string(jsb), result.Code, result.Msg)
	}

	return &result, nil
}

func DoRequestWebsocket(ctx context.Context, method, url string, headers map[string]string, params interface{}) error {

	jsb, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("json.Marshal(%+v) failed: %s", params, err.Error())
	}
	if params == nil {
		jsb = []byte{}
	}
	log.Debugf("[SECRET] %s %s -> \n%s\n", method, url, string(jsb))

	req, err := http.NewRequest(method, url, bytes.NewReader(jsb))

	if err != nil {
		return fmt.Errorf("new requrest (%s, %s, %s) failed: %s", method, url, string(jsb), err.Error())
	}

	if ctx != nil {
		req = req.WithContext(ctx)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	client.Timeout = time.Minute

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read from resp body failed: %s", err.Error())
	}

	fmt.Println(string(bs))
	if string(bs) != "" {
		return fmt.Errorf("push error")
	}

	return nil
}
