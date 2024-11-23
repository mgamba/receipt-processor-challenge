package main

import (
  "io/ioutil"
  "net/http/httptest"
  "testing"
  "os"
  "strings"
  "encoding/json"
	"fmt"
)

func TestPostReceipt(t *testing.T) {
  fname := "./examples/morning-receipt.json"
  body, _ := os.ReadFile(fname)

  ts := httptest.NewServer(initRouter())
  defer ts.Close()
  client := ts.Client()
  res, err := client.Post(fmt.Sprintf("%v/receipts/process", ts.URL), "application/json", strings.NewReader(string(body)))
  if err != nil {
    t.Errorf("Error: %v", err)
	}
  data, err := ioutil.ReadAll(res.Body)
  defer res.Body.Close()
  if err != nil {
    t.Errorf("Error: %v", err)
  }

  expected := `{"id":"bdf833c5-a299-50ef-af8a-7642b2b3a546"}`
  if string(data) != expected {
    t.Errorf("Expected %s but got %v", expected, string(data))
  }
}

func TestGetPoints(t *testing.T) {
  fname := "./examples/morning-receipt.json"
  body, _ := os.ReadFile(fname)

  ts := httptest.NewServer(initRouter())
  defer ts.Close()
  client := ts.Client()
  postPath := fmt.Sprintf("%v/receipts/process", ts.URL)
  postResponse, err := client.Post(postPath, "application/json", strings.NewReader(string(body)))
  if err != nil {
    t.Errorf("Error: %v", err)
	}
  postResponseBody, err := ioutil.ReadAll(postResponse.Body)
  postResponse.Body.Close()
  if err != nil {
    t.Errorf("Error: %v", err)
  }
  var parsedPost map[string]interface{}
	err = json.Unmarshal([]byte(postResponseBody), &parsedPost)
	if err != nil {
    t.Errorf("Error: %v", err)
	}
  id := parsedPost["id"]

  getPath := fmt.Sprintf("%v/receipts/%s/points", ts.URL, id)
  getResponse, err := client.Get(getPath)
  if err != nil {
    t.Errorf("Error: %v", err)
	}
  getResponseBody, err := ioutil.ReadAll(getResponse.Body)
  defer getResponse.Body.Close()
  if err != nil {
    t.Errorf("Error: %v", err)
  }
  var parsedGet map[string]interface{}
	err = json.Unmarshal([]byte(getResponseBody), &parsedGet)
	if err != nil {
    t.Errorf("Error: %v", err)
	}
  points := int(parsedGet["points"].(float64))

  expected := 15
  if points != expected {
    t.Errorf("Expected %d but got %d", expected, points)
  }
}

