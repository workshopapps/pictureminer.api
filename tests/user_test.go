package test

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestSignUp(t *testing.T) {
	url := "http://localhost:9000/signup"

	payload := strings.NewReader(`{` + " " + ` "user_name":"mikey",` + " " + ` "email":"michael@gmail.com", ` + " " + ` "password": "MyPassword123"` + " " + ` }`)

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}

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

	if m["user_name"] != "mikey" {
		t.Log("exp:", "mikey")
		t.Log("got:", m["user_name"])
		t.Fatal("Usernames don't match")
	}

	if m["email"] != "michael@gmail.com" {
		t.Log("exp:", "michael@gmail.com")
		t.Log("got:", m["email"])
		t.Fatal("emails don't match")
	}

	if res.StatusCode != 200 {
		t.Log("exp:", 200)
		t.Log("got:", res.StatusCode)
		t.Fatal("status codes don't match")
	}
}
