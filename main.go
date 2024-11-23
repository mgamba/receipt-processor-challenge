package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
  "github.com/google/uuid"
  "fetch/take-home/receipt"
	"bytes"
  "encoding/json"
)

var DB = make(map[string]int)

func PostReceipt(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  var reqBodyCopy bytes.Buffer
  reqBody := io.TeeReader(req.Body, &reqBodyCopy)

  if parsedReceipt, err := receipt.Parse(reqBody); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("{}"))
  } else {
    id := sha1(reqBodyCopy.String())
    DB[id] = parsedReceipt.Points()
    response, _ := json.Marshal(map[string]string{ "id":  id })
    io.WriteString(w, string(response))
  }
}

func GetPoints(w http.ResponseWriter, req *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  id := req.PathValue("id")

  if points, ok := DB[id]; ok {
    response, _ := json.Marshal(map[string]int{ "points":  points })
    io.WriteString(w, string(response))
  } else {
    w.WriteHeader(http.StatusNotFound)
    w.Write([]byte("{}"))
  }
}

func initRouter() *http.ServeMux {
  router := http.NewServeMux()
	router.HandleFunc("POST /receipts/process", PostReceipt)
	router.HandleFunc("GET /receipts/{id}/points", GetPoints)
  return router
}

func main() {
	err := http.ListenAndServe("0.0.0.0:3333", initRouter())
  if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func sha1(data string) string {
  return uuid.NewSHA1(uuid.Nil, []byte(data)).String()
}

