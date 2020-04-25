package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Html reads the html of a page as a string
func Html(url string) ([]byte, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %w", err)
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}
	return body, nil
}
