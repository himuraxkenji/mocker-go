package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type EndpointMock struct {
	Endpoint string `json:"endpoint,omitempty"`
	Method   string `json:"method,omitempty"`
	Response struct {
		Status int `json:"status,omitempty"`
		Body   any `json:"body,omitempty"`
	} `json:"response,omitempty"`
}

func InitServer() {
	responses := readResponseFiles()

	for _, r := range responses {
		route := fmt.Sprintf("%s %s/", strings.ToUpper(r.Method), r.Endpoint)
		fmt.Println(route)
		http.HandleFunc(route, func(w http.ResponseWriter, rq *http.Request) {
			if rq.Method != r.Method {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(r.Response.Status)
			json.NewEncoder(w).Encode(r.Response.Body)
		})
	}

	http.ListenAndServe(":8000", nil)

}

func readResponseFiles() []EndpointMock {
	d := "/home/himura/code/go/tools/mocker-go/internal/responses/"

	entries, err := os.ReadDir(d)
	var result []EndpointMock

	if err != nil {
		fmt.Println("Error when read directory -", err)
		return result
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		fn := e.Name()

		if strings.HasSuffix(fn, ".json") {
			fp := filepath.Join(d, fn)
			c, err := os.ReadFile(fp)

			if err != nil {
				fmt.Println("Error on read file json", err)
				continue
			}
			var r []EndpointMock
			err = json.Unmarshal(c, &r)
			if err != nil {
				fmt.Printf("a error has occur when decode - %v\n", err)
			}
			result = append(result, r...)
		}
	}

	fmt.Println("read all files in directory - ")

	return result

}
