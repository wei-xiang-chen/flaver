package utils

import (
	"bytes"
	"encoding/json"
	"flaver/globals"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
)

func RestfulSend(method, url, body string, header map[string]string) ([]byte, int, error) {
	payload := strings.NewReader(body)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		globals.GetLogger().Errorf("[RestfulSend] NewRequest error: %v", err)
		return nil, 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		globals.GetLogger().Errorf("[RestfulSend] client.Do error: %v", err)
		return nil, res.StatusCode, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		globals.GetLogger().Errorf("[RestfulSend] read response body error: %v", err)
		return nil, res.StatusCode, err
	}

	json, _ := json.Marshal(res.Header)
	globals.GetLogger().Warnf("[RestfulSend] url: %v, header: %v", url, string(json))
	return resBody, res.StatusCode, nil
}

func PostWithFormData(method, url string, body map[string]string, header map[string]string) ([]byte, int, error) {

	payload := new(bytes.Buffer)
	w := multipart.NewWriter(payload)
	for k, v := range body {
		w.WriteField(k, v)
	}
	w.Close()
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		globals.GetLogger().Warnf("[PostWithFormData] NewRequest error: %v", err)
		return nil, req.Response.StatusCode, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	for k, v := range header {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		globals.GetLogger().Warnf("[PostWithFormData] client.Do error: %v", err)
		return nil, req.Response.StatusCode, err
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		globals.GetLogger().Warnf("[PostWithFormData] read response body error: %v", err)
		return nil, req.Response.StatusCode, err
	}
	return resBody, req.Response.StatusCode, nil
}
