package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestMineServiceUpload(t *testing.T) {
	url := "http://localhost:9000/api/v1/mine-service/url"

	method := "POST"

	payload := strings.NewReader(`{
    "url": "https://images.pexels.com/photos/1561020/pexels-photo-1561020.jpeg"
}`)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2NjkzNzcyOTYsImlkIjoiT2JqZWN0SUQoXCI2MzgwOWUwMjlmMjc3MTg1NjJiMGE3NGFcIikifQ.b1M000NRCbQWN9TxlqvvCX_5khisQHDqNSfK8Igtil4")

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		t.Error(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
		return
	}
	m := make(map[string]interface{})
	marshalErr := json.Unmarshal(body, &m)
	if marshalErr != nil {
		fmt.Println(marshalErr)
		t.Error(marshalErr)
		return
	}
	t.Log(m)
	if res.StatusCode != 200 {
		t.Log("exp:", 200)
		t.Log("got:", res.StatusCode)
		t.Fatal("status codes don't match")
	}
}
