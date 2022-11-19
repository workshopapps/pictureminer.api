package mineservice

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// write unit test
func TestUnit(t *testing.T) {

}

// write endpoint test
func TestMineServiceAPI(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "{ \"status\": \"good\"}")
	}
	req := httptest.NewRequest("POST", "/api/mine-service", nil) //nil will be the request to post to the url
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if 200 != resp.StatusCode {
		t.Fatal("status code not OK")
	}
}
